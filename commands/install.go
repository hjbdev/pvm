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
	"regexp"
	"strings"
)

func Install(args []string) {
	if len(args) < 2 {
		theme.Error("You must specify a version to install.")
		return
	}

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

	desiredVersionNumbers := common.GetVersion(args[1], desireThreadSafe, "")

	if desiredVersionNumbers == (common.Version{}) {
		theme.Error("Invalid version specified")
		return
	}

	// Get the desired version from the user input
	desiredMajorVersion := desiredVersionNumbers.Major
	desiredMinorVersion := desiredVersionNumbers.Minor
	desiredPatchVersion := desiredVersionNumbers.Patch

    func fetchAndProcess(url string) []Version {
    	// Perform GET request
    	resp, err := http.Get(url)
    	if err != nil {
    		log.Fatalln(err)
    	}
    	defer resp.Body.Close()
    
    	// Read response body
    	body, err := io.ReadAll(resp.Body)
    	if err != nil {
    		log.Fatalln(err)
    	}
    
    	// Convert body to string
    	sb := string(body)
    
    	// Regex to match links
    	re := regexp.MustCompile(`<A HREF="([a-zA-Z0-9./-]+)">([a-zA-Z0-9./-]+)</A>`)
    	matches := re.FindAllStringSubmatch(sb, -1)
    
    	var versions []Version
    
        func shouldSkip(name string) bool {
        	// Skip based on specific criteria
        	return name == "" ||
        		(len(name) > 15 && strings.HasPrefix(name, "php-devel-pack-")) ||
        		(len(name) > 15 && strings.HasPrefix(name, "php-debug-pack-")) ||
        		(len(name) > 14 && strings.HasPrefix(name, "php-test-pack-")) ||
        		strings.Contains(name, "src") ||
        		!strings.HasSuffix(name, ".zip") ||
        		!strings.Contains(name, "x64")
        }
    
    	for _, match := range matches {
    		url := match[1]
    		name := match[2]
    
    		// Apply filters
    		if shouldSkip(name) {
    			continue
    		}
    
    		threadSafe := true
    		if strings.Contains(name, "nts") || strings.Contains(name, "NTS") {
    			threadSafe = false
    		}
    
    		// Append the version after processing
    		versions = append(versions, common.GetVersion(name, threadSafe, url))
    	}
    
    	return versions
    }
    
    urls := []string{
		"https://windows.php.net/downloads/releases/archives/",
		"https://windows.php.net/downloads/releases/",
	}

	var allVersions []Version

	// Fetch and process each URL
	for _, url := range urls {
		versions := fetchAndProcess(url)
		allVersions = append(allVersions, versions...)
	}

	// find desired version
	var desiredVersion common.Version

	if desiredMajorVersion > -1 && desiredMinorVersion > -1 && desiredPatchVersion > -1 {
		desiredVersion = FindExactVersion(allVersions, desiredMajorVersion, desiredMinorVersion, desiredPatchVersion, desireThreadSafe)
	}

	if desiredMajorVersion > -1 && desiredMinorVersion > -1 && desiredPatchVersion == -1 {
		desiredVersion = FindLatestPatch(allVersions, desiredMajorVersion, desiredMinorVersion, desireThreadSafe)
	}

	if desiredMajorVersion > -1 && desiredMinorVersion == -1 && desiredPatchVersion == -1 {
		desiredVersion = FindLatestMinor(allVersions, desiredMajorVersion, desireThreadSafe)
	}

	if desiredVersion == (common.Version{}) {
		theme.Error(fmt.Sprintf("Could not find the desired version: %s %s", args[1], threadSafeString))
		return
	}

	fmt.Printf("Installing PHP %s\n", desiredVersion)

	// get current dir
	currentDir, err := os.Executable()

	if err != nil {
		log.Fatalln(err)
	}
	
	fullDir := filepath.Dir(currentDir)

	// check if .pvm folder exists
	pvmPath := filepath.Join(currentDir, ".pvm")
	if _, err := os.Stat(pvmPath); os.IsNotExist(err) {
		theme.Info("Creating .pvm folder in home directory")
		os.Mkdir(pvmPath, 0755)
	}

	// check if .pvm/versions folder exists
	versionsPath := filepath.Join(pvmPath, "versions")
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		theme.Info("Creating .pvm/versions folder in home directory")
		os.Mkdir(versionsPath, 0755)
	}

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
