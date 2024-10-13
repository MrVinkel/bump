package bump

import (
	"sort"
)

func GetLatestVersion(repo *Repo) (*Version, error) {
	tags, err := repo.GetTags()
	if err != nil {
		return nil, err
	}

	Debug("tags: %v\n", PrintSlice(tags))

	versions := make([]Version, 0)
	for _, t := range tags {
		v, err := ParseVersion(t)
		if err != nil {
			Error("invalid tag: %s\n", t)
			continue
		}
		versions = append(versions, *v)
	}

	Debug("parsed versions: %s\n", PrintVersionSlice(versions))

	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) > 0
	})

	Debug("sorted versions: %v\n", PrintVersionSlice(versions))

	version := NewVersion(0, 0, 0)
	if len(versions) > 0 {
		version = &versions[0]
	}
	return version, nil
}
