package common

import (
	"errors"
	"fmt"
	"io"
	"log"
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

func RetrievePHPVersions() ([]Version, error) {
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
