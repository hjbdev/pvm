package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func (v Version) String() string {
	semantic := fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
	if v.ThreadSafe {
		return semantic + " thread safe"
	}
	return semantic + " non-thread safe"
}

func GetVersion(text string, safe bool, url string) Version {
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

func GetCurrentVersionFolder() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	content, err := os.ReadFile(filepath.Join(homeDir, ".pvm", "version"))

	if err != nil {
		return ""
	}

	return string(content)
}

func ReadPhpIni(path string) string {
	file, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(file)
}

type ExtensionStatus int

const (
	ExtensionEnabled  ExtensionStatus = 1
	ExtensionDisabled ExtensionStatus = 2
	ExtensionNotFound ExtensionStatus = 3
)

func GetExtensionStatus(ini string, extension string) (ExtensionStatus, int) {
	lines := regexp.MustCompile(`\r?\n`).Split(ini, -1)

	for index, line := range lines {
		extensionMatches := regexp.MustCompile(`extension\s*=\s*["']?([^"']+)["']?`).FindStringSubmatch(line)

		if len(extensionMatches) == 0 {
			continue
		}

		if extensionMatches[1] == extension {
			fullLine := lines[index]
			noWhitespace := strings.TrimSpace(fullLine)

			if strings.HasPrefix(noWhitespace, ";") {
				return ExtensionDisabled, index
			} else {
				return ExtensionEnabled, index
			}
		}
	}

	return ExtensionNotFound, -1
}
