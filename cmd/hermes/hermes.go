package main

import (
	"fmt"
	"os"

	"github.com/ElecTwix/hermes/pkg/gitmanager"
)

func main() {
	args := os.Args

	fmt.Println("Hello, Hermes!", args)

	workspace, ok := os.LookupEnv("GITHUB_WORKSPACE")
	if !ok {
		fmt.Println("GITHUB_WORKSPACE not set")
		os.Exit(1)
	}

	fmt.Println("Workspace: ", workspace)

	gitManager := gitmanager.NewGitManager()
	repo, err := gitManager.PlainOpen(workspace)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error opening git repository")
		os.Exit(1)
	}

	sha, ok := os.LookupEnv("GITHUB_SHA")
	if !ok {
		fmt.Println("GITHUB_SHA not set")
		os.Exit(1)
	}

	fmt.Println("SHA: ", sha)

	commitMsg, err := repo.GetCommit(sha)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting commit message")
		os.Exit(1)
	}

	fmt.Println(commitMsg)
}
