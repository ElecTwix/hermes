package gitmng

import "github.com/go-git/go-git/v5"

type GitMng struct{}

func NewGitMng() *GitMng {
	return &GitMng{}
}

func (g *GitMng) PlainOpen(path string) error {
	git.PlainOpen(path)
}

func (g *GitMng) Clone() error {
	return nil
}

func (g *GitMng) Pull() error {
	return nil
}

func (g *GitMng) Push() error {
	return nil
}
