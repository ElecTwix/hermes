package gitmanager

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitManager struct{}

func NewGitManager() *GitManager {
	return &GitManager{}
}

func (g *GitManager) PlainOpen(path string) (*GitRepository, error) {
	fmt.Println(path + "/.git")
	repo, err := git.PlainOpen(path + "/.git")
	if err != nil {
		return nil, err
	}

	return &GitRepository{Repository: repo}, nil
}

func (g *GitManager) Clone() error {
	return nil
}

func (g *GitManager) Pull() error {
	return nil
}

func (g *GitManager) Push() error {
	return nil
}

type GitRepository struct {
	Repository *git.Repository
}

func (g *GitRepository) GetCommit(sha string) string {
	hash := plumbing.NewHash(sha)
	commit, err := g.Repository.CommitObject(hash)
	if err != nil {
		return ""
	}

	fmt.Println(commit)
	return commit.String()
}
