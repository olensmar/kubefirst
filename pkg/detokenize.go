package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Detokenize - Translate tokens by values on a given path
func DetokenizeV2(path string) {

	err := filepath.Walk(path, DetokenizeDirectoryV2)
	if err != nil {
		log.Panic(err)
	}
}

// DetokenizeDirectory - Translate tokens by values on a directory level.
func DetokenizeDirectoryV2(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	if strings.Contains(path, ".gitClient") || strings.Contains(path, ".terraform") || strings.Contains(path, ".git/") {
		return nil
	}

	matched, err := filepath.Match("*", fi.Name())

	if err != nil {
		log.Panic(err)
	}

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}

		// config := configs.ReadConfig()

		newContents := string(read)

		//! registry
		gitopsGitUrl := viper.GetString("github.repo.gitops.giturl")
		awsRegion := viper.GetString("aws.region")
		awsAccountId := viper.GetString("aws.account-id")
		awsHostedZoneName := viper.GetString("aws.hosted-zone-name")
		adminEmail := viper.GetString("admin-email")
		chartmuseumIngressUrl := fmt.Sprintf("https://chartmuseum.%s", awsHostedZoneName)
		argocdIngressUrl := fmt.Sprintf("https://argocd.%s", awsHostedZoneName)
		argocdIngressUrlNoHttps := fmt.Sprintf("argocd.%s", awsHostedZoneName)
		gitopsUrlNoHttps := fmt.Sprintf("github.com/%s/gitops.git", viper.GetString("github.owner"))
		argoIngressUrl := fmt.Sprintf("https://argo.%s", awsHostedZoneName)
		vaultIngressUrl := fmt.Sprintf("https://vault.%s", awsHostedZoneName)
		vouchIngressUrl := fmt.Sprintf("https://vouch.%s", awsHostedZoneName)
		kubefirstIngressUrl := fmt.Sprintf("kubefirst.%s", awsHostedZoneName)
		// myVar := viper.GetString("aws.iam-arn")
		// myVar := viper.GetString("aws.hosted-zone-id")
		// myVar := viper.GetString("aws.hosted-zone-name")
		// myVar := viper.GetString("aws.profile")
		// myVar := viper.GetString("argocd.local.service")
		// myVar := viper.GetString("cloud-provider")
		// myVar := viper.GetString("git-provider")

		// myVar := viper.GetString("github.atlantis.webhook.secret")
		// myVar := viper.GetString("github.repo.gitops.url")
		// myVar := viper.GetString("github.repo.metaphor.url")
		// myVar := viper.GetString("github.repo.metaphor-frontend.url")
		// myVar := viper.GetString("github.repo.metaphor-go.url")
		// myVar := viper.GetString("github.owner")
		// myVar := viper.GetString("github.user")

		// //! hack
		// // viper.Set("gitops-template.repo.branch", gitopsTemplateBranchFlag)
		// viper.Set("template-repo.gitops.branch", "domain-refactor")
		// viper.Set("template-repo.gitops.url", gitopsTemplateUrlFlag)
		// // todo accommodate metaphor branch and repo override more intelligently
		// viper.Set("template-repo.metaphor.url", fmt.Sprintf("https://github.com/%s/metaphor.git", "kubefirst"))
		// viper.Set("template-repo.metaphor.branch", "main")
		// viper.Set("template-repo.metaphor-frontend.url", fmt.Sprintf("https://github.com/%s/metaphor-frontend.git", "kubefirst"))
		// viper.Set("template-repo.metaphor-frontend.branch", "main")
		// viper.Set("template-repo.metaphor-go.url", fmt.Sprintf("https://github.com/%s/metaphor-go.git", "kubefirst"))
		// viper.Set("template-repo.metaphor-go.branch", "main")
		// viper.Set("kubefirst.bot.password", kbotPassword)
		// viper.Set("kubefirst.bot.private-key", sshPrivateKey)
		// viper.Set("kubefirst.bot.public-key", sshPublicKey)
		// viper.Set("kubefirst.bot.user", "kbot")
		// viper.Set("kubefirst.state-store.bucket", k1StateStoreBucketName)
		// viper.Set("kubefirst.telemetry", useTelemetryFlag)
		// viper.Set("cluster-name", clusterNameFlag)
		// viper.Set("vault.local.service", config.LocalVaultURL)
		//! registry
		newContents = strings.Replace(newContents, "<FULL_GITOPS_REPO_GIT_URL>", gitopsGitUrl, -1)
		newContents = strings.Replace(newContents, "<FULL_GITOPS_REPO_URL_NO_HTTPS>", gitopsUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<AWS_DEFAULT_REGION>", awsRegion, -1)
		newContents = strings.Replace(newContents, "<AWS_ACCOUNT_ID>", awsAccountId, -1)
		newContents = strings.Replace(newContents, "<EMAIL_ADDRESS>", adminEmail, -1)
		newContents = strings.Replace(newContents, "<CHARTMUSEUM_INGRESS_URL>", chartmuseumIngressUrl, -1)
		newContents = strings.Replace(newContents, "<ARGOCD_INGRESS_URL_NO_HTTPS>", argocdIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<KUBEFIRST_INGRESS_URL_NO_HTTPS>", kubefirstIngressUrl, -1)

		//! terraform
		// ? argocd ingress url might be in registry?
		newContents = strings.Replace(newContents, "<ARGOCD_INGRESS_URL>", argocdIngressUrl, -1)

		// didnt see
		newContents = strings.Replace(newContents, "<ARGO_WORKFLOWS_INGRESS_URL_NO_HTTPS>", argoIngressUrl, -1)
		newContents = strings.Replace(newContents, "<VAULT_INGRESS_URL>", vaultIngressUrl, -1)
		newContents = strings.Replace(newContents, "<VOUCH_INGRESS_URL>", vouchIngressUrl, -1)
		// <ARGOCD_INGRESS_URL>
		// <ARGO_WORKFLOWS_INGRESS_URL_NO_HTTPS>
		// <GITHUB_HOST>
		// <GITHUB_USER>
		// <AWS_DEFAULT_REGION>
		// <GITHUB_OWNER>
		// <AWS_ACCOUNT_ID>
		// <AWS_DEFAULT_REGION>
	}

	return nil
}
