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
	semver        = `^(?P<prefix>0|[^0-9]*)(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	prefix        = 1
	major         = 2
	minor         = 3
	patch         = 4
	prerelease    = 5
	buildmetadata = 6
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
	fmt.Printf("test: %v\n", re.SubexpNames())
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid version format")
	}
	fmt.Printf("matches: %d\n", len(matches))

	for i, name := range re.SubexpNames() {
		fmt.Printf("match %d: %s = %s\n", i, name, matches[i])
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
		return 1
	} else if v1.Major < v2.Major {
		return -1
	}

	if v1.Minor > v2.Minor {
		return 1
	} else if v1.Minor < v2.Minor {
		return -1
	}

	if v1.Patch > v2.Patch {
		return 1
	} else if v1.Patch < v2.Patch {
		return -1
	}

	return 0
}

func (v *Version) String() string {
	if v.Prefix == nil {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%s%d.%d.%d", *v.Prefix, v.Major, v.Minor, v.Patch)
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
