package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ElecTwix/hermes/pkg/genai"
	"github.com/ElecTwix/hermes/pkg/github"
)

func main() {
	genaiInstance := genai.NewGenAI()
	geminiToken := os.Getenv("INPUT_GEMINI_TOKEN")
	if geminiToken == "" {
		fmt.Println("GEMINI_TOKEN not set")
		os.Exit(1)
	}

	ctx := context.Background()
	err := genaiInstance.Login(ctx, geminiToken, "gemini-1.5-pro")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error logging in to GenAI")
		os.Exit(1)
	}

	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		fmt.Println("GITHUB_WORKSPACE not set")
		os.Exit(1)
	}

	fmt.Println("Workspace: ", workspace)

	/*
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
	*/

	// Create the commenter
	token := os.Getenv("GH_TOKEN")
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
	changes, err := client.GetPRChanges(repoOwner, repoName, PRNumber)
	if err != nil {
		log.Fatal(err)
	}

	const promptPrefix string = `Please summarize this git commit message on markdown for my PR 
                like this and don't put the placeholder text down below:
                File: path/to/file:15-20 added http server for serving static files
                File: path/to/another/file:5-10 fixed bug with http server not serving files
                File: path/to/third/file:30-35 added new feature for support bulk create for DB. 
        `

	strChanges := ""
	for i, change := range changes {
		fmt.Printf("Change %d: %s\n", i, change.GetCommit().GetMessage())
		for _, file := range change.Files {
			strChanges += fmt.Sprintf("File: %s:%d-%d %s\n", file.GetFilename(), file.GetAdditions(), file.GetDeletions(), file.GetStatus())
		}

	}

	prompt := fmt.Sprintf("%s DATA: [%s]", promptPrefix, strChanges)
	modelOutput, err := genaiInstance.Generate(prompt, ctx)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error generating content")
		os.Exit(1)
	}

	err = client.CommentOnPR(repoOwner, repoName, PRNumber, modelOutput)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error commenting on PR")
		os.Exit(1)
	}

	fmt.Println(modelOutput)

	fmt.Println("Commented on PR successfully")
}
