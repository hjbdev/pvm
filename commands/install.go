package commands

import (
	"archive/zip"
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Install(args []string) {
	if len(args) < 2 {
		theme.Error("You must specify a version to install.")
		return
	}

	version := args[1]

	pvmPath, versionsPath, err := setupDirectories()
	if err != nil {
		log.Fatalf("Failed to set up directories: %v", err)
	}

	// check if version is already installed
	versionPath := filepath.Join(versionsPath, version)
	if _, err := os.Stat(versionPath); err == nil {
		theme.Error(fmt.Sprintf("PHP %s is already installed.", version))
		return
	}

	if runtime.GOOS == "darwin" {
		installForMac(version, pvmPath, versionsPath)
	} else if runtime.GOOS == "windows" {
		installForWindows(args, pvmPath, versionsPath)
	} else {
		theme.Error(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
	}
}

func setupDirectories() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	pvmPath := filepath.Join(homeDir, ".pvm")
	if _, err := os.Stat(pvmPath); os.IsNotExist(err) {
		theme.Info("Creating .pvm folder in home directory")
		if err := os.Mkdir(pvmPath, 0755); err != nil {
			return "", "", err
		}
	}

	versionsPath := filepath.Join(pvmPath, "versions")
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		theme.Info("Creating .pvm/versions folder")
		if err := os.Mkdir(versionsPath, 0755); err != nil {
			return "", "", err
		}
	}

	return pvmPath, versionsPath, nil
}

func installForMac(version string, pvmPath string, versionsPath string) {
	theme.Info(fmt.Sprintf("Installing PHP %s for macOS using Homebrew...", version))

	// Check if brew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		theme.Error("Homebrew is not installed. Please install it to continue.")
		theme.Info("See: https://brew.sh/")
		return
	}

	// Tap shivammathur/php for older versions
	versionParts := strings.Split(version, ".")
	if len(versionParts) >= 2 {
		major := versionParts[0]
		minor := versionParts[1]
		// Example logic: tap for versions older than 8.0
		if major < "8" || (major == "7" && minor <= "4") {
			theme.Info("Tapping shivammathur/php for older PHP versions...")
			cmd := exec.Command("brew", "tap", "shivammathur/php")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				theme.Error("Failed to tap shivammathur/php.")
				log.Printf("Error: %v\n", err)
				return
			}
		}
	}

	phpPackage := fmt.Sprintf("php@%s", version)
	theme.Info(fmt.Sprintf("Running 'brew install %s'...", phpPackage))
	cmd := exec.Command("brew", "install", phpPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		theme.Error(fmt.Sprintf("Failed to install PHP %s using Homebrew.", version))
		log.Printf("Error: %v\n", err)
		return
	}

	// Get brew prefix
	prefixCmd := exec.Command("brew", "--prefix")
	prefixBytes, err := prefixCmd.Output()
	if err != nil {
		theme.Error("Failed to get brew prefix.")
		log.Printf("Error: %v\n", err)
		return
	}
	brewPrefix := strings.TrimSpace(string(prefixBytes))
	phpInstallPath := filepath.Join(brewPrefix, "opt", phpPackage)

	// Create symlinks in .pvm/versions
	versionDir := filepath.Join(versionsPath, version)
	if err := os.Mkdir(versionDir, 0755); err != nil {
		theme.Error(fmt.Sprintf("Failed to create directory for PHP %s.", version))
		log.Printf("Error: %v\n", err)
		return
	}

	binPath := filepath.Join(phpInstallPath, "bin")
	filesToLink := []string{"php", "phpize", "php-config", "pecl"}
	for _, file := range filesToLink {
		source := filepath.Join(binPath, file)
		destination := filepath.Join(versionDir, file)
		if err := os.Symlink(source, destination); err != nil {
			// Don't fail if a file doesn't exist (e.g. older versions might not have all of them)
			theme.Warning(fmt.Sprintf("Could not create symlink for %s: %v", file, err))
		}
	}

	theme.Success(fmt.Sprintf("Finished installing PHP %s", version))
	theme.Info(fmt.Sprintf("Run 'pvm use %s' to start using it.", version))
}

func installForWindows(args []string, pvmPath string, versionsPath string) {
	desireThreadSafe := true
	if len(args) > 2 {
		if args[2] == "nts" {
			desireThreadSafe = false
		}
	}

	var threadSafeString string
	if desireThreadSafe {
		threadSafeString = "thread safe"
	} else {
		threadSafeString = "non-thread safe"
	}

	if desireThreadSafe {
		theme.Warning("Thread safe version will be installed")
	} else {
		theme.Warning("Non-thread safe version will be installed")
	}

	desiredVersionNumbers := common.ComputeVersion(args[1], desireThreadSafe, "")

	if desiredVersionNumbers == (common.Version{}) {
		theme.Error("Invalid version specified")
		return
	}

	// Get the desired version from the user input
	desiredMajorVersion := desiredVersionNumbers.Major
	desiredMinorVersion := desiredVersionNumbers.Minor
	desiredPatchVersion := desiredVersionNumbers.Patch

	versions, err := common.RetrievePHPVersions()
	if err != nil {
		log.Fatalln(err)
	}

	// find desired version
	var desiredVersion common.Version

	if desiredMajorVersion > -1 && desiredMinorVersion > -1 && desiredPatchVersion > -1 {
		desiredVersion = FindExactVersion(versions, desiredMajorVersion, desiredMinorVersion, desiredPatchVersion, desireThreadSafe)
	}

	if desiredMajorVersion > -1 && desiredMinorVersion > -1 && desiredPatchVersion == -1 {
		desiredVersion = FindLatestPatch(versions, desiredMajorVersion, desiredMinorVersion, desireThreadSafe)
	}

	if desiredMajorVersion > -1 && desiredMinorVersion == -1 && desiredPatchVersion == -1 {
		desiredVersion = FindLatestMinor(versions, desiredMajorVersion, desireThreadSafe)
	}

	if desiredVersion == (common.Version{}) {
		theme.Error(fmt.Sprintf("Could not find the desired version: %s %s", args[1], threadSafeString))
		return
	}

	fmt.Printf("Installing PHP %s\n", desiredVersion)

	theme.Info("Downloading")

	// zip filename from url
	zipUrl := "https://windows.php.net" + desiredVersion.Url
	zipFileName := strings.Split(desiredVersion.Url, "/")[len(strings.Split(desiredVersion.Url, "/"))-1]
	zipPath := filepath.Join(versionsPath, zipFileName)

	// check if zip already exists
	if _, err := os.Stat(zipPath); err == nil {
		theme.Error(fmt.Sprintf("PHP %s already exists", desiredVersion))
		return
	}

	// Get the data
	if _, err := downloadFile(zipUrl, zipPath); err != nil {
		log.Fatalf("Error while downloading PHP from %v: %v!", zipUrl, err)
	}

	// extract the zip file to a folder
	phpFolder := strings.Replace(zipFileName, ".zip", "", -1)
	phpPath := filepath.Join(versionsPath, phpFolder)
	theme.Info("Unzipping")
	Unzip(zipPath, phpPath)

	// remove the zip file
	theme.Info("Cleaning up")
	err = os.Remove(zipPath)
	if err != nil {
		log.Fatalln(err)
	}

	// install composer
	composerFolderPath := filepath.Join(phpPath, "composer")
	if _, err := os.Stat(composerFolderPath); os.IsNotExist(err) {
		theme.Info("Creating composer folder")
		os.Mkdir(composerFolderPath, 0755)
	}

	composerPath := filepath.Join(composerFolderPath, "composer.phar")
	composerUrl := "https://getcomposer.org/download/latest-stable/composer.phar"
	if desiredVersion.LessThan(common.Version{Major: 7, Minor: 2}) {
		composerUrl = "https://getcomposer.org/download/latest-2.2.x/composer.phar"
	}

	if _, err := downloadFile(composerUrl, composerPath); err != nil {
		log.Fatalf("Error while downloading Composer from %v: %v!", composerUrl, err)
	}

	theme.Success(fmt.Sprintf("Finished installing PHP %s", desiredVersion))
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindExactVersion(versions []common.Version, major int, minor int, patch int, threadSafe bool) common.Version {
	for _, version := range versions {
		if version.ThreadSafe != threadSafe {
			continue
		}
		if version.Major == major && version.Minor == minor && version.Patch == patch {
			return version
		}
	}

	return common.Version{}
}

func FindLatestPatch(versions []common.Version, major int, minor int, threadSafe bool) common.Version {
	latestPatch := common.Version{}

	for _, version := range versions {
		if version.ThreadSafe != threadSafe {
			continue
		}
		if version.Major == major && version.Minor == minor {
			if latestPatch.Patch == -1 || version.Patch > latestPatch.Patch {
				latestPatch = version
			}
		}
	}

	return latestPatch
}

func FindLatestMinor(versions []common.Version, major int, threadSafe bool) common.Version {
	latestMinor := common.Version{}

	for _, version := range versions {
		if version.ThreadSafe != threadSafe {
			continue
		}
		if version.Major == major {
			if latestMinor.Minor == -1 || version.Minor > latestMinor.Minor {
				if latestMinor.Patch == -1 || version.Patch > latestMinor.Patch {
					latestMinor = version
				}
			}
		}
	}

	return latestMinor
}

func downloadFile(fileUrl string, filePath string) (bool, error) {
	downloadResponse, err := http.Get(fileUrl)
	if err != nil {
		return false, err
	}

	defer downloadResponse.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return false, err
	}

	// Write the body to file
	_, err = io.Copy(out, downloadResponse.Body)

	if err != nil {
		out.Close()
		return false, err
	}

	// Close the file
	out.Close()
	return true, nil
}
