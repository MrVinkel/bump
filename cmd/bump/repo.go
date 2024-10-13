package bump

import (
	"fmt"

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
	tagRefs.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, r.Name().Short())
		return nil
	})
	return tags, nil
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
