package internal

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
	prefixMajorMinorPatch = `^([^0-9]+)*([0-9]+)\.([0-9]+)\.([0-9]+)$`
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
	re := regexp.MustCompile(prefixMajorMinorPatch)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid version format")
	}

	i := 1
	var prefix *string
	if len(matches) == 5 {
		prefix = Ptr(matches[i])
		i++
	}

	major, err := strconv.Atoi(matches[i])
	if err != nil {
		return nil, err
	}

	i++
	minor, err := strconv.Atoi(matches[i])
	if err != nil {
		return nil, err
	}

	i++
	patch, err := strconv.Atoi(matches[i])
	if err != nil {
		return nil, err
	}

	return &Version{
		Prefix: prefix,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
	}, nil
}

func BumpPatch(v *Version) *Version {
	return NewVersionP(v.Prefix, v.Major, v.Minor, v.Patch+1)
}

func BumpMinor(v *Version) *Version {
	return NewVersionP(v.Prefix, v.Major, v.Minor+1, 0)
}

func BumpMajor(v *Version) *Version {
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
