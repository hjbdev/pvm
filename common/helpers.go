package common

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	Url        string
	ThreadSafe bool
}

func (v Version) Semantic() string {
	if v.Patch == -1 {
		return fmt.Sprintf("%v.%v", v.Major, v.Minor)
	}
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}

func (v Version) StringShort() string {
	semantic := v.Semantic()
	if v.ThreadSafe {
		return semantic
	}
	return semantic + " nts"
}

func (v Version) String() string {
	semantic := v.Semantic()
	if v.ThreadSafe {
		return semantic + " thread safe"
	}
	return semantic + " non-thread safe"
}

func ComputeVersion(text string, safe bool, url string) Version {
	versionRe := regexp.MustCompile(`([0-9]{1,3})(?:.([0-9]{1,3}))?(?:.([0-9]{1,3}))?`)
	matches := versionRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return Version{}
	}

	major, err := strconv.Atoi(matches[0][1])
	if err != nil {
		major = -1
	}

	minor, err := strconv.Atoi(matches[0][2])
	if err != nil {
		minor = -1
	}

	patch, err := strconv.Atoi(matches[0][3])
	if err != nil {
		patch = -1 // Default patch to -1 if not present
	}

	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		ThreadSafe: safe,
		Url:        url,
	}
}

func (v Version) Compare(o Version) int {
	if v.Major == -1 || o.Major == -1 {
		return 0
	}
	if v.Major != o.Major {
		if v.Major < o.Major {
			return -1
		}
		return 1
	}

	if v.Minor == -1 || o.Minor == -1 {
		return 0
	}
	if v.Minor != o.Minor {
		if v.Minor < o.Minor {
			return -1
		}
		return 1
	}

	if v.Patch == -1 || o.Patch == -1 {
		return 0
	}
	if v.Patch != o.Patch {
		if v.Patch < o.Patch {
			return -1
		}
		return 1
	}

	return 0
}

func (v Version) CompareThreadSafe(o Version) int {
	result := v.Compare(o)
	if result != 0 {
		return result
	}

	if v.ThreadSafe == o.ThreadSafe {
		return 0
	}

	if v.ThreadSafe {
		return -1
	}
	return 1
}

func (v Version) LessThan(o Version) bool {
	return v.CompareThreadSafe(o) == -1
}

func (v Version) Same(o Version) bool {
	return v.CompareThreadSafe(o) == 0
}

func SortVersions(input []Version) []Version {
	sort.SliceStable(input, func(i, j int) bool {
		return input[i].LessThan(input[j])
	})
	return input
}

func RetrievePHPVersions() ([]Version, error) {
	if runtime.GOOS == "darwin" {
		return retrieveMacPHPVersions()
	}
	return retrieveWindowsPHPVersions()
}

func retrieveMacPHPVersions() ([]Version, error) {
	// Search in both homebrew/core and shivammathur/php
	cmd1 := exec.Command("brew", "search", "--formulae", "php")
	out1, err1 := cmd1.CombinedOutput()
	// It's okay if one of them fails (e.g. tap not installed), as long as we get some versions.
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to search homebrew/core for php versions: %v\n", err1)
	}

	cmd2 := exec.Command("brew", "search", "--formulae", "shivammathur/php/php")
	out2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to search shivammathur/php for php versions: %v\n", err2)
	}

	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("failed to execute 'brew search' in both homebrew/core and shivammathur/php")
	}

	output := string(out1) + "\n" + string(out2)
	lines := strings.Split(output, "\n")
	versions := make([]Version, 0)
	seen := make(map[string]bool)

	// Matches php@8.1, shivammathur/php/php@8.1 and extracts the version number
	re := regexp.MustCompile(`php@(\d+\.\d+)`)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// We only care about formulae, which are not surrounded by {}
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "{") {
			continue
		}

		matches := re.FindStringSubmatch(trimmedLine)
		var versionString string
		if len(matches) > 1 {
			versionString = matches[1]
		} else if strings.HasSuffix(trimmedLine, "php") {
			// Handle the case for unversioned formulae like "php" or "shivammathur/php/php", which point to the latest version.
			infoCmd := exec.Command("brew", "info", trimmedLine, "--json=v1")
			infoOut, infoErr := infoCmd.Output()
			if infoErr == nil {
				// A bit of a hacky json parse to find the version
				jsonRe := regexp.MustCompile(`"version"\s*:\s*"([0-9.]+)"`)
				jsonMatches := jsonRe.FindStringSubmatch(string(infoOut))
				if len(jsonMatches) > 1 {
					fullVersion := jsonMatches[1]
					parts := strings.Split(fullVersion, ".")
					if len(parts) >= 2 {
						versionString = fmt.Sprintf("%s.%s", parts[0], parts[1])
					}
				}
			}
			if versionString == "" {
				continue // could not determine version
			}
		} else {
			continue
		}

		if _, exists := seen[versionString]; !exists {
			versions = append(versions, ComputeVersion(versionString, true, ""))
			seen[versionString] = true
		}
	}

	return versions, nil
}

func retrieveWindowsPHPVersions() ([]Version, error) {
	// perform get request to https://windows.php.net/downloads/releases/archives/
	resp, err := http.Get("https://windows.php.net/downloads/releases/archives/")
	if err != nil {
		return nil, err
	}
	// We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Convert the body to type string
	sb := string(body)

	// regex match
	re := regexp.MustCompile(`<A HREF="([a-zA-Z0-9./-]+)">([a-zA-Z0-9./-]+)</A>`)
	matches := re.FindAllStringSubmatch(sb, -1)

	versions := make([]Version, 0)

	for _, match := range matches {
		url := match[1]
		name := match[2]

		// check if name starts with "php-devel-pack-"
		if name != "" && len(name) > 15 && name[:15] == "php-devel-pack-" {
			continue
		}
		// check if name starts with "php-debug-pack-"
		if name != "" && len(name) > 15 && name[:15] == "php-debug-pack-" {
			continue
		}
		// check if name starts with "php-test-pack-"
		if name != "" && len(name) > 15 && name[:14] == "php-test-pack-" {
			continue
		}

		// check if name contains "src"
		if name != "" && strings.Contains(name, "src") {
			continue
		}

		// check if name does not end in zip
		if name != "" && !strings.HasSuffix(name, ".zip") {
			continue
		}

		threadSafe := true

		// check if name contains "nts" or "NTS"
		if name != "" && (strings.Contains(name, "nts") || strings.Contains(name, "NTS")) {
			threadSafe = false
		}

		// make sure we only get x64 versions
		if name != "" && !strings.Contains(name, "x64") {
			continue
		}

		// regex match name and push to versions
		versions = append(versions, ComputeVersion(name, threadSafe, url))
	}
	return versions, nil
}

func RetrieveInstalledPHPVersions() ([]Version, error) {
	versions := make([]Version, 0)
	// get users home dir
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
		return versions, err
	}

	// check if .pvm folder exists
	pvmPath := filepath.Join(homeDir, ".pvm")
	if _, err := os.Stat(pvmPath); os.IsNotExist(err) {
		return versions, errors.New("no PHP versions installed")
	}

	// check if .pvm/versions folder exists
	versionsPath := filepath.Join(pvmPath, "versions")
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		return versions, errors.New("no PHP versions installed")
	}

	// get all folders in .pvm/versions
	folders, err := os.ReadDir(versionsPath)
	if err != nil {
		return versions, err
	}

	for _, folder := range folders {
		folderName := folder.Name()
		safe := true
		if strings.Contains(folderName, "nts") || strings.Contains(folderName, "NTS") {
			safe = false
		}

		versions = append(versions, ComputeVersion(folderName, safe, ""))
	}
	SortVersions(versions)
	return versions, nil
}
