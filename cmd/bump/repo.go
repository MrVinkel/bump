package bump

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type Repo struct {
	repo *git.Repository
}

func NewRepo(path string) (*Repo, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	return &Repo{repo: repo}, nil
}

func (r *Repo) GetTags() ([]string, error) {
	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}

	tags := []string{}
	err = tagRefs.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, r.Name().Short())
		return nil
	})
	return tags, err
}

func (r *Repo) CreateTag(tag string) error {
	head, err := r.repo.Head()
	if err != nil {
		return err
	}
	_, err = r.repo.CreateTag(tag, head.Hash(), nil)
	return err
}

func (r *Repo) PushTag(tag string) error {
	refSpec := fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag)
	return r.repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(refSpec)},
	})
}

func (r *Repo) CreateAndPushTag(tag string) error {
	err := r.CreateTag(tag)
	if err != nil {
		return err
	}
	return r.PushTag(tag)
}

func (r *Repo) Fetch() error {
	err := r.repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func (r *Repo) HasChanges() (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := w.Status()
	if err != nil {
		return false, err
	}
	return !status.IsClean(), nil
}

func (r *Repo) IsSynced() (bool, error) {
	head, err := r.repo.Head()
	if err != nil {
		return false, err
	}
	name := strings.TrimPrefix(head.Name().String(), "refs/heads/")

	refs, err := r.repo.References()
	if err != nil {
		return false, err
	}
	// find remote branch
	var remote *plumbing.Reference
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsRemote() {
			return nil
		}
		if strings.TrimPrefix(ref.Name().String(), "refs/remotes/origin/") == name {
			remote = ref
		}
		return nil
	})

	if remote == nil {
		return false, fmt.Errorf("remote branch for %s not found", name)
	}
	if remote.Hash() != head.Hash() {
		return false, nil
	}

	return true, err
}
