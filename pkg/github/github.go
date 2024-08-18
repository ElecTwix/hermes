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

	name := "Hermes Summary Bot"

	// Create the comment object
	comment := &github.IssueComment{
		Body: &commentContent,
		User: &github.User{
			Name: &name,
		},
	}

	// Create the comment on the pull request
	_, _, err := g.Client.Issues.CreateComment(ctx, repoOwner, repoName, prNumber, comment)
	if err != nil {
		return err
	}

	return nil
}

func (g *GithubClient) GetPRChanges(repoOwner, repoName string, prNumber int) ([]*github.RepositoryCommit, error) {
	// Create GitHub client
	ctx := context.Background()
	// Get the commits in the pull request
	commits, _, err := g.Client.PullRequests.ListCommits(ctx, repoOwner, repoName, prNumber, nil)
	if err != nil {
		return nil, err
	}
	return commits, nil
}
