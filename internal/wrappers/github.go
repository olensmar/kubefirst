package wrappers

import (
	"errors"
	"log"
	"os"

	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/handlers"
)

// AuthenticateGitHubUserWrapper receives a handler that was previously instantiated, and communicate with GitHub.
// This wrapper is necessary to avoid code repetition when requesting GitHub PAT or Access token.
func AuthenticateGitHubUserWrapper(config *configs.Config, gitHubHandler *handlers.GitHubHandler) (string, error) {

	githubToken := config.GithubToken
	if githubToken != "" {
		return githubToken, nil
	}

	githubToken, err := gitHubHandler.AuthenticateUser()
	if err != nil {
		return "", err
	}

	if githubToken == "" {
		return "", errors.New("unable to retrieve a GitHub token for the user")
	}

	if err := os.Setenv("GITHUB_TOKEN", githubToken); err != nil {
		return "", err
	}
	log.Println("\nGITHUB_TOKEN set via OAuth")

	return githubToken, nil
}
