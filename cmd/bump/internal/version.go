package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Prefix     *string
	Major      int
	Minor      int
	Patch      int
	PreRelease []string
	Build      *string
}

const (
	// semver regex from semver.org plus a prefix group
	semver = `^(?P<prefix>0|[^0-9]*)(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	// match groups
	prefix        = 1
	major         = 2
	minor         = 3
	patch         = 4
	prerelease    = 5
	buildmetadata = 6
	// compare
	greator = 1
	less    = -1
	equal   = 0
)

func NewVersion(prefix *string, major, minor, patch int, preRelease []string, build *string) *Version {
	return &Version{
		Prefix:     prefix,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Build:      build,
	}
}

func ParseVersion(version string) (*Version, error) {
	re := regexp.MustCompile(semver)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid version format")
	}

	prefixStr := matches[prefix]
	var prefix *string
	if prefixStr != "" {
		prefix = &prefixStr
	}

	major, err := strconv.Atoi(matches[major])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(matches[minor])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(matches[patch])
	if err != nil {
		return nil, err
	}

	var preRelease []string
	if matches[prerelease] != "" {
		preRelease = strings.Split(matches[prerelease], ".")
	}

	buildStr := matches[buildmetadata]
	var build *string
	if buildStr != "" {
		build = &buildStr
	}

	return &Version{
		Prefix:     prefix,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Build:      build,
	}, nil
}

func BumpPatch(v *Version) *Version {
	return NewVersion(v.Prefix, v.Major, v.Minor, v.Patch+1, []string{}, v.Build)
}

func BumpMinor(v *Version) *Version {
	return NewVersion(v.Prefix, v.Major, v.Minor+1, 0, []string{}, v.Build)
}

func BumpMajor(v *Version) *Version {
	return NewVersion(v.Prefix, v.Major+1, 0, 0, []string{}, v.Build)
}

func Compare(v1, v2 Version) int {
	if v1.Major > v2.Major {
		return greator
	} else if v1.Major < v2.Major {
		return less
	}

	if v1.Minor > v2.Minor {
		return greator
	} else if v1.Minor < v2.Minor {
		return less
	}

	if v1.Patch > v2.Patch {
		return greator
	} else if v1.Patch < v2.Patch {
		return less
	}

	v1Len := len(v1.PreRelease)
	v2Len := len(v2.PreRelease)
	if v1Len > v2Len {
		return less
	}

	for i := range v2Len {
		if i >= v1Len {
			if v1Len == 0 {
				return greator
			}
			return less
		}

		v1Num, v1Ok := tryParseNumber(v1.PreRelease[i])
		v2Num, v2Ok := tryParseNumber(v2.PreRelease[i])

		if v1Ok && v2Ok {
			if v1Num > v2Num {
				return greator
			} else if v1Num < v2Num {
				return less
			}
		} else if !v1Ok && v2Ok {
			return greator
		} else if v1Ok && !v2Ok {
			return less
		} else if v1.PreRelease[i] > v2.PreRelease[i] {
			return greator
		} else if v1.PreRelease[i] < v2.PreRelease[i] {
			return less
		}
	}

	return equal
}

func tryParseNumber(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (v *Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prefix != nil {
		version = fmt.Sprintf("%s%s", *v.Prefix, version)
	}
	if len(v.PreRelease) > 0 {
		version = fmt.Sprintf("%s-%s", version, strings.Join(v.PreRelease, "."))
	}
	if v.Build != nil {
		version = fmt.Sprintf("%s+%s", version, *v.Build)
	}
	return version
}

func VersionSliceString(s []Version) string {
	b := strings.Builder{}
	first := false
	for _, v := range s {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}
		b.WriteString(v.String())
	}
	return b.String()
}
