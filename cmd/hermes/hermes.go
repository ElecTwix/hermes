package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ElecTwix/hermes/pkg/genai"
	"github.com/ElecTwix/hermes/pkg/github"
	"github.com/ElecTwix/hermes/pkg/gitmanager"
)

func main() {
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
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

	sha := os.Getenv("PR_GITHUB_SHA")
	if sha == "" {
		fmt.Println("PR_GITHUB_SHA not set")
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
	apiKey := os.Getenv("API_KEY_SECRET")
	if apiKey == "" {
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
	modelOutput, err := genaiInstance.Generate(prompt, ctx)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error generating content")
		os.Exit(1)
	}

	fmt.Println(modelOutput)

	// Create the commenter
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN not set")
		os.Exit(1)
	}

	PRNumberStr := os.Getenv("GITHUB_PR_NUMBER")
	if PRNumberStr == "" {
		fmt.Println("GITHUB_PR_NUMBER not set")
		os.Exit(1)
	}

	repoOwner := os.Getenv("GITHUB_REPO_OWNER")
	if repoOwner == "" {
		fmt.Println("GITHUB_REPO_OWNER not set")
		os.Exit(1)
	}

	repoName := os.Getenv("GITHUB_REPO_NAME")
	if repoName == "" {
		fmt.Println("GITHUB_REPO_NAME not set")
		os.Exit(1)
	}

	PRNumber, err := strconv.Atoi(PRNumberStr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error converting PR number to int")
		os.Exit(1)
	}

	client := github.NewGithubClient()
	client.Auth(token)
	err = client.CommentOnPR(repoOwner, repoName, PRNumber, modelOutput)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error commenting on PR")
		os.Exit(1)
	}

	fmt.Println("Commented on PR successfully")
}
