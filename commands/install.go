package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type Version struct {
	Major      string
	Minor      string
	Patch      string
	Url        string
	ThreadSafe bool
}

func Install(args []string) {
	if len(args) < 2 {
		color.Red("You must specify a version to install.")
		return
	}

	versionRe := regexp.MustCompile(`([0-9]{1,3})(?:.([0-9]{1,3}))?(?:.([0-9]{1,3}))?`)

	desiredVersionMatches := versionRe.FindAllStringSubmatch(args[1], -1)

	if len(desiredVersionMatches) == 0 {
		color.Red("Invalid version specified")
		return
	}

	// Get the desired version from the user input
	desiredMajorVersion := desiredVersionMatches[0][1]
	desiredMinorVersion := desiredVersionMatches[0][2]
	desiredPatchVersion := desiredVersionMatches[0][3]

	// perform get request to https://windows.php.net/downloads/releases/archives/
	resp, err := http.Get("https://windows.php.net/downloads/releases/archives/")
	if err != nil {
		log.Fatalln(err)
	}
	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
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

		// regex match name
		versionMatches := versionRe.FindAllStringSubmatch(name, -1)

		major := versionMatches[0][1]
		minor := versionMatches[0][2]
		patch := versionMatches[0][3]

		// push to versions
		versions = append(versions, Version{
			Major:      major,
			Minor:      minor,
			Patch:      patch,
			Url:        url,
			ThreadSafe: threadSafe,
		})
	}

	// find desired version
	desiredVersion := Version{}

	for _, version := range versions {
		// Exact match
		if version.Major == desiredMajorVersion && version.Minor == desiredMinorVersion && version.Patch == desiredPatchVersion {
			fmt.Println("Exact match found")
			desiredVersion = version
			break
		}
		// Major and minor version match, find the highest patch version
		if version.Major == desiredMajorVersion && version.Minor == desiredMinorVersion {
			if version.Patch > desiredVersion.Patch {
				desiredVersion = version
			}
		}
		// Major version matches, find the highest patch version for a 0 minor version
		if version.Major == desiredMajorVersion {
			if version.Minor == "0" {
				desiredVersion = version
			}
			if version.Minor == desiredMinorVersion {
				if version.Patch > desiredVersion.Patch {
					desiredVersion = version
				}
			}
		}
	}

	fmt.Println(desiredVersion)
}
