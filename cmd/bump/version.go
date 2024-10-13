package bump

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Prefix *string
	Major  int
	Minor  int
	Patch  int
}

const (
	majorMinorPatch       = `^([0-9]+)\.([0-9]+)\.([0-9]+)$`
	prefixMajorMinorPatch = `^([^0-9]+)([0-9]+)\.([0-9]+)\.([0-9]+)$`
)

func NewVersion(major, minor, patch int) *Version {
	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func NewVersionP(prefix *string, major, minor, patch int) *Version {
	return &Version{
		Prefix: prefix,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
	}
}

func ParseVersion(version string) (*Version, error) {
	if version[0] >= '0' && version[0] <= '9' {
		return parse(version)
	}
	return parsePrefix(version)
}

func parse(version string) (*Version, error) {
	re := regexp.MustCompile(majorMinorPatch)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid version format")
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func parsePrefix(version string) (*Version, error) {
	re := regexp.MustCompile(prefixMajorMinorPatch)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid version format")
	}

	prefix := matches[1]

	major, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, err
	}

	return &Version{
		Prefix: &prefix,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
	}, nil
}

func (v *Version) BumpPatch() *Version {
	return NewVersionP(v.Prefix, v.Major, v.Minor, v.Patch+1)
}

func (v *Version) BumpMinor() *Version {
	return NewVersionP(v.Prefix, v.Major, v.Minor+1, 0)
}

func (v *Version) BumpMajor() *Version {
	return NewVersionP(v.Prefix, v.Major+1, 0, 0)
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
