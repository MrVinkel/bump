package bump

import (
	"errors"
	"sort"
)

func GetLatestVersion(repo *Repo) (*Version, error) {
	tags, err := repo.GetTags()
	if err != nil {
		return nil, err
	}

	Debug("tags: %v\n", SliceString(tags))

	versions := make([]Version, 0)
	for _, t := range tags {
		v, err := ParseVersion(t)
		if err != nil {
			Error("invalid tag: %s\n", t)
			Debug("error: %v\n", err)
			continue
		}
		versions = append(versions, *v)
	}

	Debug("parsed versions: %s\n", VersionSliceString(versions))

	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) > 0
	})

	Debug("sorted versions: %v\n", VersionSliceString(versions))

	version := NewVersion(0, 0, 0)
	if len(versions) > 0 {
		version = &versions[0]
	}
	return version, nil
}

func CheckRepositoryStatus(repo *Repo) error {
	if *NoVerify {
		return nil
	}

	hasChanages, err := repo.HasChanges()
	if err != nil {
		return err
	}
	if hasChanages {
		return errors.New("uncommitted changes")
	}

	if !*NoFetch {
		if err = repo.Fetch(); err != nil {
			return err
		}
	}

	synced, err := repo.IsSynced()
	if err != nil {
		return err
	}
	if !synced {
		return errors.New("unpushed changes")
	}
	return nil
}
