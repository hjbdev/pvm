package common

import (
	"fmt"
	"regexp"
	"strconv"
	"os"
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
	Folder os.DirEntry
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

func (v Version) LessThan(other Version) bool {
	return v.Compare(other) == -1
}

func (v Version) LessThanOrEqual(other Version) bool {
	return v.Compare(other) == -1 || v.Compare(other) == 0
}

func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) == 1
}

func (v Version) GreaterThanOrEqual(other Version) bool {
	return v.Compare(other) == 1 || v.Compare(other) == 0
}

func (v Version) Equal(other Version) bool {
	return v.Compare(other) == 0
}

func (v Version) Same(other Version) bool {
	return v.Compare(other) == 0 && v.ThreadSafe == other.ThreadSafe
}
