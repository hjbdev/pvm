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
	if len(args) < 1 {
		theme.Error("You must specify a version to install.")
		return
	}

	desireThreadSafe := true
	installPath := "" // This will store the install path

	// Check for "nts" flag or installation path
	if len(args) > 1 {
		if args[1] == "nts" {
			desireThreadSafe = false
		} else {
			// If it's not "nts", consider it as the installation path
			installPath = args[1]
		}
	}
	if len(args) > 2 {
		if args[2] == "nts" {
			desireThreadSafe = false
		} else {
			// If it's not "nts", consider it as the installation path
			installPath = args[2]
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

	desiredVersionNumbers := common.ComputeVersion(args[0], desireThreadSafe, "")

	if desiredVersionNumbers == (common.Version{}) {
		theme.Error("Invalid version specified")
		return
	}

	// Get the desired version from the user input
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
		theme.Error(fmt.Sprintf("Could not find the desired version: %s %s", args[0], threadSafeString))
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
