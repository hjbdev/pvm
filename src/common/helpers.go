package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

type VersionMeta struct {
	Number Version
	Folder string
}

type VersionJson struct {
	Path    string
	Source  string
	Version string
}

func (v Version) Semantic() string {
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
		return Version{
			Major:      -1,
			Minor:      -1,
			Patch:      -1,
			ThreadSafe: safe,
			Url:        url,
		}
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
		patch = -1
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

func RetrievePHPVersions(url string) ([]Version, error) {
	// Make the HTTP request to the provided URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PHP versions from URL '%s': %v", url, err)
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from URL '%s': %v", url, err)
	}

	// Convert the body to a string
	sb := string(body)

	// Regex match to extract href links
	re := regexp.MustCompile(`<A HREF="([a-zA-Z0-9./-]+)">([a-zA-Z0-9./-]+)</A>`)
	matches := re.FindAllStringSubmatch(sb, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no valid version links found on the page '%s'", url)
	}

	// Initialize a slice to store the versions
	versions := make([]Version, 0)

	// Loop through the matches and filter out unwanted versions
	for _, match := range matches {
		versionURL := match[1]
		name := match[2]

		// Check if name starts with any of the unwanted prefixes
		if strings.HasPrefix(name, "php-devel-pack-") ||
			strings.HasPrefix(name, "php-debug-pack-") ||
			strings.HasPrefix(name, "php-test-pack-") {
			continue
		}

		// Check if the name contains "src", which is not a PHP version
		if strings.Contains(name, "src") {
			continue
		}

		// Check if the name does not end with ".zip", which indicates it's not a valid PHP version
		if !strings.HasSuffix(name, ".zip") {
			continue
		}

		// Default to thread-safe version
		threadSafe := true

		// Check if the name contains "nts" or "NTS", which indicates a non-thread-safe version
		if strings.Contains(strings.ToLower(name), "nts") {
			threadSafe = false
		}

		// Ensure we only process x64 versions
		if !strings.Contains(name, "x64") {
			continue
		}

		// Create a version object based on the extracted name and thread-safety flag
		versions = append(versions, ComputeVersion(name, threadSafe, versionURL))
	}

	// Return the list of versions
	if len(versions) == 0 {
		return nil, fmt.Errorf("no valid PHP versions found on the page '%s'", url)
	}

	return versions, nil
}

// RetrieveInstalledPHPVersions reads the installed PHP versions from the versions.json file inside .pvm/versions.
func RetrieveInstalledPHPVersions() ([]Version, error) {
	versions := make([]Version, 0)

	// Get the current dir
	currentDir, err := os.Executable()
	if err != nil {
		return versions, fmt.Errorf("unable to get current executable directory: %v", err)
	}

	fullDir := filepath.Dir(currentDir)

	// Construct the path to the versions.json file inside .pvm/versions
	jsonPath := filepath.Join(fullDir, "versions", "versions.json")

	// Check if the versions.json file exists
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		return versions, fmt.Errorf("no PHP versions installed (versions.json not found)")
	}

	// Read the versions.json file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return versions, fmt.Errorf("failed to read versions.json: %v", err)
	}

	// Parse the existing JSON data from the versions.json file
	var existingVersions []map[string]interface{}
	if err := json.Unmarshal(data, &existingVersions); err != nil {
		return versions, fmt.Errorf("failed to parse versions.json: %v", err)
	}

	// Process the existing versions and map them to Version structs
	for _, entry := range existingVersions {
		// Extract the relevant fields from the JSON entry
		versionStr, ok := entry["version"].(string)
		if !ok {
			continue
		}

		// Determine if the version is thread-safe or non-thread-safe
		safe := true
		if strings.Contains(strings.ToLower(versionStr), "nts") {
			safe = false
		}

		// Create the Version object
		version := ComputeVersion(versionStr, safe, "")

		// Add the version to the list
		versions = append(versions, version)
	}

	// Sort the versions in ascending order
	SortVersions(versions)

	return versions, nil
}

// AddToVersionsJson adds a new version entry to versions.json.
func AddToVersionsJson(path string, version string, source string) error {
	// Get the current dir
	currentDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to get current executable directory: %v", err)
	}
	fullDir := filepath.Dir(currentDir)

	// Check if .pvm/versions folder exists, if not create it
	versionsPath := filepath.Join(fullDir, "versions")
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		if err := os.Mkdir(versionsPath, 0755); err != nil {
			return fmt.Errorf("failed to create versions directory: %v", err)
		}
	}

	// Construct the path to versions.json
	jsonPath := filepath.Join(versionsPath, "versions.json")

	// Read the current versions.json (if it exists)
	var existingVersions []map[string]interface{}
	if _, err := os.Stat(jsonPath); err == nil {
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("failed to read versions.json: %v", err)
		}

		// Parse existing JSON data
		if err := json.Unmarshal(data, &existingVersions); err != nil {
			return fmt.Errorf("failed to parse versions.json: %v", err)
		}
	}

	// Check if the path already exists in the versions list
	for _, entry := range existingVersions {
		if entry["path"] == path {
			return fmt.Errorf("the PHP version from the path %s already exists in versions.json", path)
		}
	}

	// Prepare the new version entry
	versionEntry := map[string]interface{}{
		"path":    path,
		"version": version,
		"source":  source,
	}

	// Append the new version entry to existing versions
	existingVersions = append(existingVersions, versionEntry)

	// Marshal the updated data back to JSON
	updatedData, err := json.MarshalIndent(existingVersions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %v", err)
	}

	// Write the updated JSON back to versions.json
	if err := os.WriteFile(jsonPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write to versions.json: %v", err)
	}

	return nil
}

// removeFromVersionJson removes the uninstalled version from the versions.json file.
func RemoveFromVersionJson(path string) error {
	// Get current dir
	currentDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to get current executable directory: %v", err)
	}
	fullDir := filepath.Dir(currentDir)

	// Path to the versions.json file
	jsonPath := filepath.Join(fullDir, "versions", "versions.json")

	// Read the existing versions from versions.json
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read versions.json: %v", err)
	}

	// Parse the existing versions data into a slice of maps
	var existingVersions []map[string]interface{}
	if err := json.Unmarshal(data, &existingVersions); err != nil {
		return fmt.Errorf("failed to parse versions.json: %v", err)
	}

	// Check if the path exists in the versions list
	var entryFound bool
	for _, entry := range existingVersions {
		if entry["path"] == path {
			entryFound = true
			break
		}
	}

	// If the path is not found, return an error
	if !entryFound {
		return fmt.Errorf("the PHP version from the path %s was not found in versions.json", path)
	}

	// Find and remove the entry for the specified path
	for i, entry := range existingVersions {
		if entry["path"] == path {
			// Remove the entry from the slice
			existingVersions = append(existingVersions[:i], existingVersions[i+1:]...)
			break
		}
	}

	// Marshal the updated data back to JSON
	updatedData, err := json.MarshalIndent(existingVersions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated versions data: %v", err)
	}

	// Write the updated JSON data back to the versions.json file
	if err := os.WriteFile(jsonPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated versions.json: %v", err)
	}

	return nil
}
