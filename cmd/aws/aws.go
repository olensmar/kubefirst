package aws

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/downloadManager"
	"github.com/kubefirst/kubefirst/internal/reports"
	"github.com/kubefirst/kubefirst/internal/wrappers"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runAws(cmd *cobra.Command, args []string) error {

	config := configs.ReadConfig()

	//* confirm with user to continue
	// var userInput string
	// printConfirmationScreen()
	// todo skip to confirmation value if we get to here and they exit here
	// otherwise we need to clean and restart / reset
	// go counter()
	// fmt.Println("to proceed, type 'yes' any other answer will exit")
	// fmt.Scanln(&userInput)
	// if userInput != "yes" {
	// 	os.Exit(1)
	// }

	// viper.GetString("admin-email")
	// viper.GetString("aws.account-id")
	// viper.GetString("aws.iam-arn")
	// viper.GetString("aws.hosted-zone-id")
	awsHostedZone := viper.GetString("aws.hosted-zone-name")
	// viper.GetString("aws.profile")
	// viper.GetString("aws.region")
	// viper.GetString("argocd.local.service")
	// viper.GetString("cloud-provider")
	// viper.GetString("gitops-template.repo.branch")
	// viper.GetString("gitops-template.repo.url")
	// viper.GetString("git-provider")
	// viper.GetString("github.atlantis.webhook.secret")
	// githubGitopsRepoUrl := viper.GetString("github.gitops-repo.url")
	// viper.GetString("github.owner")
	// viper.GetString("github.user")
	// viper.GetString("kubefirst.bot.password")
	// viper.GetString("kubefirst.bot.private-key")
	// viper.GetString("kubefirst.bot.public-key")
	// viper.GetString("kubefirst.bot.user")
	// viper.GetString("kubefirst.state-store.bucket")
	// viper.GetString("kubefirst.telemetry")
	// viper.GetString("cluster-name")
	// viper.GetString("vault.local.service")
	// config := configs.ReadConfig()
	//* emit cluster install started
	if useTelemetryFlag {
		if err := wrappers.SendSegmentIoTelemetry(awsHostedZone, pkg.MetricMgmtClusterInstallStarted); err != nil {
			log.Println(err)
		}
	}

	// todo download dependencies
	//* download dependencies `$HOME/.k1/tools`
	if !viper.GetBool("kubefirst.dependency-download.complete") {
		log.Println("installing kubefirst dependencies")

		err := downloadManager.DownloadTools(config)
		if err != nil {
			return err
		}

		log.Println("download dependencies `$HOME/.k1/tools` complete")
		viper.Set("kubefirst.dependency-download.complete", true)
		viper.WriteConfig()
	} else {
		log.Println("download dependencies `$HOME/.k1/tools` already done - continuing")
		// log.Println("k1Config: kubefirst.dependency-download.complete") // todo is this valuable for debugging?
	}

	// todo clone and detokenize repo

	//!
	// todo terraform apply github repositories (all)

	// todo clone and detoknize repos and push to remote

	// todo terraform apply base - include additional s3 buckets for better management

	// todo detoknize kms key id and re-push local content to remote

	// todo restore ssl... also automatically backup ssl at the end

	// todo create initial argocd repository (this is the connection to argocd as a 'repo') destinationGitopsRepoUrl
	//* investigate - is this doing the same thing as pkg.CreateSSHKey where it writes a file?

	// todo helm install argocd

	// todo wait for argocd to be ready

	// todo set argocd credentials

	//! stop here before continuing
	// todo apply argocd registry

	// todo set argocd credentials
	// todo set argocd credentials

	return nil

}

// todo move below functions? pkg? rename?
func counter() {
	i := 0
	for {
		time.Sleep(time.Second * 1)
		i++
	}
}

func printConfirmationScreen() {
	var createKubefirstSummary bytes.Buffer
	createKubefirstSummary.WriteString(strings.Repeat("-", 70))
	createKubefirstSummary.WriteString("\nCreate Kubefirst Cluster?\n")
	createKubefirstSummary.WriteString(strings.Repeat("-", 70))
	createKubefirstSummary.WriteString("\n\nAWS Account Details:\n\n")
	createKubefirstSummary.WriteString(fmt.Sprintf("Account IAM Arn:  %s\n", viper.GetString("aws.iam-arn")))
	createKubefirstSummary.WriteString(fmt.Sprintf("Account ID:       %s\n", viper.GetString("aws.account-id")))
	createKubefirstSummary.WriteString(fmt.Sprintf("Hosted Zone Id:   %s\n", viper.GetString("aws.hosted-zone-id")))
	createKubefirstSummary.WriteString(fmt.Sprintf("Hosted Zone Name: %s\n", viper.GetString("aws.hosted-zone-name")))
	createKubefirstSummary.WriteString(fmt.Sprintf("Profile:          %s\n", viper.GetString("aws.profile")))
	createKubefirstSummary.WriteString(fmt.Sprintf("Region:           %s\n", viper.GetString("aws.region")))
	createKubefirstSummary.WriteString("\n\nGithub Organization Details:\n\n")
	createKubefirstSummary.WriteString(fmt.Sprintf("Organization: %s\n", viper.GetString("github.owner")))
	createKubefirstSummary.WriteString(fmt.Sprintf("User:         %s\n", viper.GetString("github.user")))
	createKubefirstSummary.WriteString("New Github Repository URL's:\n")
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("github.repo.gitops.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("github.repo.metaphor.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("github.repo.metaphor-frontend.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("github.repo.metaphor-go.url")))

	createKubefirstSummary.WriteString("\n\nTemplate Repositories URL's:\n\n")
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("gitops-template.repo.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("gitops-template.repo.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor-frontend.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor-frontend.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor-go.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor-go.branch")))

	fmt.Println(reports.StyleMessage(createKubefirstSummary.String()))
}
