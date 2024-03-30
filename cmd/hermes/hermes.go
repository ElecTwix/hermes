package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ElecTwix/hermes/pkg/genai"
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

	genaiInstance := genai.NewGenAI()
	apiKey, ok := os.LookupEnv("API_KEY_SECRET")
	if !ok {
		fmt.Println("API_KEY_SECRET not set")
		os.Exit(1)
	}

	ctx := context.Background()
	err = genaiInstance.Login(apiKey, ctx, "gemini-pro")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error logging in to GenAI")
		os.Exit(1)
	}

	const promptPrefix string = `Please summarize this git commit message for my PR: like this:
                File: path/to/file:15-20 added http server for serving static files
                File: path/to/another/file:5-10 fixed bug with http server not serving files
                File: path/to/third/file:30-35 added new feature for support bulk create for DB. 
        `

	prompt := fmt.Sprintf("%s DATA: [%s]", promptPrefix, commitMsg)
	output, err := genaiInstance.Generate(prompt, ctx)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error generating content")
		os.Exit(1)
	}

	fmt.Println(output)
}
