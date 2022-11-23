package local

import (
	"log"
	"net/http"

	"github.com/kubefirst/kubefirst/internal/handlers"
	"github.com/kubefirst/kubefirst/internal/services"
	"github.com/spf13/cobra"
)

func validateAws(cmd *cobra.Command, args []string) error {

	// config := configs.ReadConfig()

	// if err := pkg.ValidateK1Folder(config.K1FolderPath); err != nil {
	// 	return err
	// }

	gitProviderFlag, err := cmd.Flags().GetString("git-provider")
	if err != nil {
		return err
	}
	log.Println("git-provider flag value", gitProviderFlag)

	gitHubAccessToken := config.GitHubPersonalAccessToken
	httpClient := http.DefaultClient
	gitHubService := services.NewGitHubService(httpClient)
	gitHubHandler := handlers.NewGitHubHandler(gitHubService)

	return nil
}
