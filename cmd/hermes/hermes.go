package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ElecTwix/hermes/pkg/gitmanager"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	fmt.Println(commitMsg)
}
