package github

import (
	"context"

	"github.com/google/go-github/v60/github"
)

type GithubClient struct {
	Client *github.Client
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		Client: github.NewClient(nil),
	}
}

func (g *GithubClient) Auth(token string) {
	g.Client = github.NewClient(nil).WithAuthToken(token)
}

func (g *GithubClient) CommentOnPR(repoOwner, repoName string, prNumber int, commentContent string) error {
	// Create GitHub client
	ctx := context.Background()

	// Create the comment object
	comment := &github.PullRequestComment{
		Body: &commentContent,

		// Optional: For line comments
		// Path:     github.String("file.txt"),
		// Position: github.Int(10),
	}

	// Create the comment on the pull request
	_, _, err := g.Client.PullRequests.CreateComment(ctx, repoOwner, repoName, prNumber, comment)
	if err != nil {
		return err
	}

	return nil
}
