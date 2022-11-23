package aws

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/downloadManager"
	"github.com/kubefirst/kubefirst/internal/handlers"
	"github.com/kubefirst/kubefirst/internal/services"
	"github.com/kubefirst/kubefirst/internal/wrappers"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateAws(cmd *cobra.Command, args []string) error {

	// infoCmd.Run(cmd, args)

	config := configs.ReadConfig()

	adminEmailFlag, err := cmd.Flags().GetString("admin-email")
	if err != nil {
		log.Panic(err)
	}

	cloudProviderFlag, err := cmd.Flags().GetString("cloud-provider")
	if err != nil {
		log.Panic(err)
	}

	clusterNameFlag, err := cmd.Flags().GetString("cluster-name")
	if err != nil {
		log.Panic(err)
	}

	githubOwner, err := cmd.Flags().GetString("github-owner")
	if err != nil {
		log.Panic(err)
	}

	gitopsTemplateUrlFlag, err := cmd.Flags().GetString("gitops-template-url")
	if err != nil {
		log.Panic(err)
	}

	gitopsTemplateBranchFlag, err := cmd.Flags().GetString("gitops-template-branch")
	if err != nil {
		log.Panic(err)
	}
	gitProviderFlag, err := cmd.Flags().GetString("git-provider")
	if err != nil {
		log.Panic(err)
	}

	hostedZoneNameFlag, err := cmd.Flags().GetString("hosted-zone-name")
	if err != nil {
		log.Panic(err)
	}

	if strings.HasSuffix(hostedZoneNameFlag, ".") {
		hostedZoneNameFlag = hostedZoneNameFlag[:len(hostedZoneNameFlag)-1]
	}
	log.Println("hostedZoneNameFlag:", hostedZoneNameFlag)

	useTelemetryFlag, err := cmd.Flags().GetBool("use-telemetry")
	if err != nil {
		log.Panic(err)
	}

	if useTelemetryFlag {
		if err := wrappers.SendSegmentIoTelemetry("", pkg.MetricInitStarted); err != nil {
			log.Println(err)
		}
	}

	// if err := pkg.ValidateK1Folder(config.K1FolderPath); err != nil {
	// 	return err
	// }

	log.Println("git-provider flag value", gitProviderFlag)

	httpClient := http.DefaultClient
	githubToken := config.GithubToken
	gitHubService := services.NewGitHubService(httpClient)
	gitHubHandler := handlers.NewGitHubHandler(gitHubService)
	// get GitHub data to set user and owner based on the provided token
	githubUser, err := gitHubHandler.GetGitHubUser(githubToken)
	if err != nil {
		return err
	}
	// githubToken, err := wrappers.AuthenticateGitHubUserWrapper(config, gitHubHandler)
	// if err != nil {
	// 	return err
	// }

	// // get GitHub data to set user and owner based on the provided token
	// githubUser, err := gitHubHandler.GetGitHubUser(githubToken)
	// if err != nil {
	// 	return err
	// }

	// // todo need to check the token from the user to
	// // todo see if its an admin (owner in the org)
	// _, err := wrappers.CheckGithubOrganizationPermissionsWrapper(config.GithubToken, githubOwner, githubUser)
	// if err != nil {
	// 	log.Println("insufficient permissions for the authenicated user. please ensure the token is an Owner in the organization")
	// 	return err
	// }
	err = gitHubHandler.CheckGithubOrganizationPermissions(githubToken, githubOwner, githubUser)
	if err != nil {
		// is a log here valuable or duplicative?
		// log.Println(fmt.Sprintf("insufficient permissions for the authenticated user (GITHUB_TOKEN).\n please make sure the token is an `Owner` in %s", githubOwner))
		return err
	}

	silentModeMockFlag := false
	if !viper.GetBool("kubefirst.validate.k1-file.complete") {
		pkg.InformUser("writing `$HOME/.kubefirst` file content", silentModeMockFlag)

		viper.Set("admin-email", adminEmailFlag)
		viper.Set("aws.hosted-zone-name", hostedZoneNameFlag)
		viper.Set("argocd.local.service", config.LocalArgoCdURL)
		viper.Set("cloud-provider", cloudProviderFlag)
		viper.Set("gitops-template.repo.branch", gitopsTemplateBranchFlag)
		viper.Set("gitops-template.repo.url", gitopsTemplateUrlFlag)
		viper.Set("git-provider", gitProviderFlag)
		viper.Set("github.atlantis.webhook.secret", pkg.Random(20))
		viper.Set("github.gitops-repo.url", fmt.Sprintf("https://github.com/%s/gitops.git", githubOwner))
		viper.Set("github.owner", githubOwner)
		viper.Set("github.user", githubUser)
		viper.Set("kubefirst.telemetry", useTelemetryFlag)
		viper.Set("cluster-name", clusterNameFlag)
		viper.Set("vault.local.service", config.LocalVaultURL)

		err = viper.WriteConfig()
		if err != nil {
			return err
		}

		viper.Set("kubefirst.validate.k1-file.complete", true)
		viper.WriteConfig()
	} else {
		log.Println("creation of `$HOME/.kubefirst` file already done - continuing")
		log.Println("k1Config: kubefirst.validate.k1-file.complete")
	}
	log.Println("validation and `kubefirst cli` environment is ready")

	log.Println("installing kubefirst dependencies")
	err = downloadManager.DownloadTools(config)
	if err != nil {
		return err
	}
	log.Println("dependency installation complete")

	return nil
}
