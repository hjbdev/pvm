package commands

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hjbdev/pvm/common"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// fetchLatestRelease is injectable for tests.
var fetchLatestRelease = func(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/hjbdev/pvm/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "pvm-version-check")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status %d from github api", resp.StatusCode)
	}

	var out struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.TagName == "" {
		return "", errors.New("empty tag_name in release")
	}
	return out.TagName, nil
}

// fetchLatestReleaseInfo returns the latest tag and the browser_download_url for pvm.exe asset.
var fetchLatestReleaseInfo = func(ctx context.Context) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/hjbdev/pvm/releases/latest", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "pvm-version-check")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("unexpected status %d from github api", resp.StatusCode)
	}

	var out struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	if out.TagName == "" {
		return "", "", errors.New("empty tag_name in release")
	}
	for _, a := range out.Assets {
		if strings.EqualFold(a.Name, "pvm.exe") {
			return out.TagName, a.BrowserDownloadURL, nil
		}
	}
	return out.TagName, "", errors.New("pvm.exe asset not found in release")
}

// CheckForUpdate compares current (e.g. "dev" or "v1.2.3") with latest release tag.
// Returns the latest tag (as returned by GitHub), whether it's newer than current, and error.
func CheckForUpdate(current string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tag, err := fetchLatestRelease(ctx)
	if err != nil {
		return "", false, err
	}

	cur := strings.TrimPrefix(strings.TrimSpace(current), "v")
	latest := strings.TrimPrefix(strings.TrimSpace(tag), "v")
	if cur == "dev" {
		// assume dev is older than any release
		return tag, true, nil
	}

	if compareSemver(latest, cur) > 0 {
		return tag, true, nil
	}
	return tag, false, nil
}

// InstallLatest runs the remote install script. If auto==false, it will prompt the user
// (reads from STDIN). If running on non-windows, returns an error noting unsupported OS.
func InstallLatest(auto bool) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("auto-install only supported on Windows (GOOS=%s)", runtime.GOOS)
	}
	// Allow overriding the installer source with PVM_INSTALL_SCRIPT. If the value
	// is an http(s) URL, it will be executed via `irm <url> | iex`. If it's a
	// filesystem path to a .ps1 script, run it with `-File`.
	installer := os.Getenv("PVM_INSTALL_SCRIPT")
	var cmd *exec.Cmd
	if installer == "" {
		cmdStr := "irm https://pvm.hjb.dev/install.ps1 | iex"
		// Build PowerShell command
		cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", cmdStr)
	} else {
		installer = strings.TrimSpace(installer)
		if strings.HasPrefix(installer, "http://") || strings.HasPrefix(installer, "https://") {
			cmdStr := "irm " + installer + " | iex"
			cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", cmdStr)
		} else {
			// treat as file path
			if _, err := os.Stat(installer); err != nil {
				return fmt.Errorf("installer script not found: %s", installer)
			}
			cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", installer)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if !auto {
		// ask for confirmation
		fmt.Printf("Run remote installer now? (y/N): ")
		var resp string
		if _, err := fmt.Fscanln(os.Stdin, &resp); err != nil {
			return err
		}
		r := strings.ToLower(strings.TrimSpace(resp))
		if r != "y" && r != "yes" {
			return errors.New("user declined install")
		}
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installer failed: %w", err)
	}

	// Post-install verification for script-based installer (best-effort):
	paths, _ := common.NewPVMPaths()
	dest := filepath.Join(paths.BinDir, "pvm.exe")
	if _, statErr := os.Stat(dest); statErr == nil {
		newVer, _ := getPVMVersionFromBinary(dest)
		if newVer == "" {
			// attempt restore from backup if available
			bak := dest + ".bak"
			if _, bErr := os.Stat(bak); bErr == nil {
				_ = os.Rename(bak, dest)
			}
			return fmt.Errorf("installer succeeded but verification failed: new binary returned empty version")
		}
	}

	return nil
}

// installLatestRunner is an injectable wrapper used by the CLI. Default points
// to DownloadAndInstallLatest which performs a safe download+replace flow.
var installLatestRunner = DownloadAndInstallLatest

// DownloadAndInstallLatest downloads the latest pvm.exe release asset, verifies
// an optional checksum (PVM_INSTALL_CHECKSUM), backs up the current binary and
// replaces it atomically.
func DownloadAndInstallLatest(auto bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tag, assetURL, err := fetchLatestReleaseInfo(ctx)
	if err != nil {
		return err
	}

	paths, err := common.NewPVMPaths()
	if err != nil {
		return err
	}
	destDir := paths.BinDir
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("pvm-%d.tmp", time.Now().UnixNano()))
	if assetURL == "" {
		return errors.New("no pvm.exe asset URL available")
	}

	if err := downloadFile(assetURL, tmpFile); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// optional checksum verification
	expected := os.Getenv("PVM_INSTALL_CHECKSUM")
	if expected != "" {
		f, err := os.Open(tmpFile)
		if err != nil {
			return err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		got := hex.EncodeToString(h.Sum(nil))
		if !strings.EqualFold(got, strings.TrimSpace(expected)) {
			return fmt.Errorf("checksum mismatch: got %s expected %s", got, expected)
		}
	}

	dest := filepath.Join(destDir, "pvm.exe")
	// capture and backup existing binary (if present)
	if _, err := os.Stat(dest); err == nil {
		bak := dest + ".bak"
		if err := os.Rename(dest, bak); err != nil {
			return fmt.Errorf("could not backup existing pvm.exe: %w", err)
		}
	}

	// move tmp -> dest
	if err := os.Rename(tmpFile, dest); err != nil {
		// try to remove dest and retry
		_ = os.Remove(dest)
		if err2 := os.Rename(tmpFile, dest); err2 != nil {
			// attempt restore from backup
			bak := dest + ".bak"
			if _, statErr := os.Stat(bak); statErr == nil {
				_ = os.Rename(bak, dest)
			}
			return fmt.Errorf("failed to move new binary into place: %w", err2)
		}
	}

	// Post-install verification: run the new binary and check version matches expected tag.
	newVer, _ := getPVMVersionFromBinary(dest)
	// prefer matching full tag (with or without v)
	want := strings.TrimPrefix(tag, "v")
	if newVer == "" {
		// verification failed: attempt rollback
		bak := dest + ".bak"
		if _, statErr := os.Stat(bak); statErr == nil {
			_ = os.Rename(bak, dest)
		}
		return fmt.Errorf("verification failed: new binary returned empty version")
	}
	if !strings.Contains(newVer, tag) && !strings.Contains(newVer, want) {
		// mismatch: rollback
		bak := dest + ".bak"
		if _, statErr := os.Stat(bak); statErr == nil {
			_ = os.Rename(bak, dest)
		}
		return fmt.Errorf("verification failed: installed version '%s' does not match expected '%s'", newVer, tag)
	}

	// verification passed; remove backup
	_ = os.Remove(dest + ".bak")
	return nil
}

// getPVMVersionFromBinary runs the pvm binary with 'help' and parses the Version line.
func getPVMVersionFromBinary(path string) (string, error) {
	out, err := exec.Command(path, "help").CombinedOutput()
	if err != nil {
		return "", err
	}
	s := string(out)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Version ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Version ")), nil
		}
	}
	return "", nil
}

// compareSemver compares two dot-separated numeric versions. Returns 1 if a>b, 0 if equal, -1 if a<b.
func compareSemver(a, b string) int {
	if a == b {
		return 0
	}
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		ai := 0
		bi := 0
		if i < len(as) {
			fmt.Sscanf(as[i], "%d", &ai)
		}
		if i < len(bs) {
			fmt.Sscanf(bs[i], "%d", &bi)
		}
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	return 0
}
