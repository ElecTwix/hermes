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

func (g *GitRepository) GetCommit(sha string) (string, error) {
	hash := plumbing.NewHash(sha)
	commit, err := g.Repository.CommitObject(hash)
	if err != nil {
		return "", err
	}

	// Get the parent commit of the current commit
	parentCommit, err := commit.Parent(0)
	if err != nil {
		panic(err)
	}

	// Get the trees for current and parent commits
	commitTree, err := commit.Tree()
	if err != nil {
		panic(err)
	}
	parentTree, err := parentCommit.Tree()
	if err != nil {
		panic(err)
	}

	changes, err := commitTree.Diff(parentTree)
	if err != nil {
		panic(err)
	}

	for _, change := range changes {
		from, to, err := change.Files()
		if err != nil {
			return "", err
		}

		before := from.Name
		after := to.Name
		fmt.Printf("Change: %s -> %s\n", before, after)

		// Print the changes in the files
		patch, err := change.Patch()
		if err != nil {
			return "", err
		}

		fmt.Println(patch)
	}

	// Print the changes
	fmt.Println(commit)
	return commit.Message, nil
}
