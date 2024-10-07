package common

import (
	"fmt"
	"regexp"
	"strconv"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)
}

func GetVersion(text string) Version {
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
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
