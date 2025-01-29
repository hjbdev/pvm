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
	"path/filepath"
	"strings"
)

func Install(args []string) {
	desireThreadSafe := true
	installPath := "" // This will store the install path
	var requestedVersion string

	if len(args) > 0 {
		requestedVersion = args[0]
	} else {
		requestedVersion = ""
		theme.Warning("Latest version will be installed")
	}

	if len(args) > 1 {
		// Process additional arguments (install path or "nts" flag)
		for _, arg := range args[1:] {
			if arg == "nts" {
				desireThreadSafe = false
			} else {
				installPath = arg
			}
		}
	}

	// Print the selected thread safety mode
	var threadSafeString string
	if desireThreadSafe {
		threadSafeString = "thread safe"
		theme.Warning("Thread safe version will be installed")
	} else {
		threadSafeString = "non-thread safe"
		theme.Warning("Non-thread safe version will be installed")
	}

	desiredVersionNumbers := common.ComputeVersion(requestedVersion, desireThreadSafe, "")
	// Get the desired version components
	desiredMajorVersion := desiredVersionNumbers.Major
	desiredMinorVersion := desiredVersionNumbers.Minor
	desiredPatchVersion := desiredVersionNumbers.Patch

	latestVersions, err := common.RetrievePHPVersions("https://windows.php.net/downloads/releases/")
	if err != nil {
		log.Fatalln(err)
	}
	archivesVersions, err := common.RetrievePHPVersions("https://windows.php.net/downloads/releases/archives/")
	if err != nil {
		log.Fatalln(err)
	}

	versions := append(latestVersions, archivesVersions...)

	// find desired version
	desiredVersion := FindVersion(versions, desiredMajorVersion, desiredMinorVersion, desiredPatchVersion, desireThreadSafe)

	if desiredVersion == (common.Version{}) {
		theme.Error(fmt.Sprintf("Could not find the desired version: %s %s", requestedVersion, threadSafeString))
		return
	}

	theme.Title(fmt.Sprintf("Installing PHP %s", desiredVersion))

	// Get current dir
	currentDir, err := os.Executable()

	if err != nil {
		log.Fatalln(err)
	}

	fullDir := filepath.Dir(currentDir)

	var versionsPath string
	// If no path is provided, default to current directory
	if installPath == "" {
		// Check if the install path folder exists
		versionsPath = filepath.Join(fullDir, "versions")
		if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
			theme.Info("Creating versions folder in the specified path")
			err := os.Mkdir(versionsPath, 0755)
			if err != nil {
				theme.Error(fmt.Sprintf("Failed to create the versions directory: %v", err))
				return
			}
		}
	} else {
		versionsPath = installPath
	}

	theme.Info("Downloading")

	// Zip filename from URL
	zipUrl := "https://windows.php.net" + desiredVersion.Url
	zipFileName := strings.Split(desiredVersion.Url, "/")[len(strings.Split(desiredVersion.Url, "/"))-1]
	zipPath := filepath.Join(versionsPath, zipFileName)

	// Check if zip already exists
	if _, err := os.Stat(zipPath); err == nil {
		theme.Error(fmt.Sprintf("PHP %s already exists", desiredVersion))
		return
	}

	// Get the data
	if _, err := downloadFile(zipUrl, zipPath); err != nil {
		theme.Error(fmt.Sprintf("Error while downloading PHP from %v: %v!", zipUrl, err))
	}

	// Extract the zip file to a folder
	phpFolder := strings.Replace(zipFileName, ".zip", "", -1)
	phpPath := filepath.Join(versionsPath, phpFolder)
	theme.Info("Unzipping")
	Unzip(zipPath, phpPath)

	// Remove the zip file
	theme.Info("Cleaning up")
	err = os.Remove(zipPath)
	if err != nil {
		theme.Error("Error while cleaning up zip file")
	}

	// Install composer
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
		theme.Error(fmt.Sprintf("Error while downloading Composer from %v: %v!", composerUrl, err))
	}

	// Add to versions.json
	common.AddToVersionsJson(phpPath, desiredVersion.String(), "pvm")

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

func FindVersion(versions []common.Version, major, minor, patch int, threadSafe bool) common.Version {
	var latestMajor, latestMinor, latestPatch int = -1, -1, -1

	// Case 1: All are -1 → Find the latest available version (highest major, minor, patch)
	if major == -1 && minor == -1 && patch == -1 {
		for _, version := range versions {
			if version.ThreadSafe != threadSafe {
				continue
			}
			if version.Major > latestMajor ||
				(version.Major == latestMajor && version.Minor > latestMinor) ||
				(version.Major == latestMajor && version.Minor == latestMinor && version.Patch > latestPatch) {
				latestMajor, latestMinor, latestPatch = version.Major, version.Minor, version.Patch
			}
		}
		major, minor, patch = latestMajor, latestMinor, latestPatch
	}

	// Case 2: Minor and Patch are -1 → Find the latest minor and patch for the given major
	if minor == -1 && patch == -1 {
		for _, version := range versions {
			if version.ThreadSafe != threadSafe || version.Major != major {
				continue
			}
			if version.Minor > latestMinor ||
				(version.Minor == latestMinor && version.Patch > latestPatch) {
				latestMinor, latestPatch = version.Minor, version.Patch
			}
		}
		minor, patch = latestMinor, latestPatch
	}

	// Case 3: Patch is -1 → Find the latest patch for the given major and minor
	if patch == -1 {
		for _, version := range versions {
			if version.ThreadSafe != threadSafe || version.Major != major || version.Minor != minor {
				continue
			}
			if version.Patch > latestPatch {
				latestPatch = version.Patch
			}
		}
		patch = latestPatch
	}

	// Case 4: All values are provided → Look for an exact match
	for _, version := range versions {
		if version.ThreadSafe == threadSafe &&
			version.Major == major &&
			version.Minor == minor &&
			version.Patch == patch {
			return version
		}
	}

	// No matching version found
	return common.Version{}
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
