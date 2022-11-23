package aws

import (
	"errors"
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

	githubOwnerFlag, err := cmd.Flags().GetString("github-owner")
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
		if err := wrappers.SendSegmentIoTelemetry(hostedZoneNameFlag, pkg.MetricInitStarted); err != nil {
			log.Println(err)
		}
	}

	// if err := pkg.ValidateK1Folder(config.K1FolderPath); err != nil {
	// 	return err
	// }

	httpClient := http.DefaultClient
	githubToken := config.GithubToken
	if len(githubToken) == 0 {
		return errors.New("ephemeral tokens not supported for cloud installations, please set a GITHUB_TOKEN environment variable to continue\n https://docs.kubefirst.io/kubefirst/github/install.html#step-3-kubefirst-init")
	}
	gitHubService := services.NewGitHubService(httpClient)
	gitHubHandler := handlers.NewGitHubHandler(gitHubService)
	// get GitHub data to set user based on the provided token
	log.Println("verifying github user")
	githubUser, err := gitHubHandler.GetGitHubUser(githubToken)
	if err != nil {
		return err
	}
	log.Println("github user is: ", githubUser)

	err = gitHubHandler.CheckGithubOrganizationPermissions(githubToken, githubOwnerFlag, githubUser)
	if err != nil {
		// is a log here valuable or duplicative?
		// log.Println(fmt.Sprintf("insufficient permissions for the authenticated user (GITHUB_TOKEN).\n please make sure the token is an `Owner` in %s", githubOwnerFlag))
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
		viper.Set("github.gitops-repo.url", fmt.Sprintf("https://github.com/%s/gitops.git", githubOwnerFlag))
		viper.Set("github.owner", githubOwnerFlag)
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
