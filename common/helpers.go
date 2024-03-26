package common

import (
	"fmt"
	"regexp"
)

type Version struct {
	Major string
	Minor string
	Patch string
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

	return Version{
		Major: matches[0][1],
		Minor: matches[0][2],
		Patch: matches[0][3],
	}
}
