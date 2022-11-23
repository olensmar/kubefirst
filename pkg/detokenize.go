package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubefirst/kubefirst/configs"
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

		config := configs.ReadConfig()

		newContents := string(read)

		//! registry
		//! from viper
		gitopsGitUrl := viper.GetString("github.repo.gitops.giturl")
		awsRegion := viper.GetString("aws.region")
		awsAccountId := viper.GetString("aws.account-id")
		awsHostedZoneName := viper.GetString("aws.hosted-zone-name")
		atlantisWebhookUrl := viper.GetString("github.atlantis.webhook.url")
		adminEmail := viper.GetString("admin-email")
		chartmuseumStorageBucketName := "chartmuseum-bucket-name" // todo blocking
		clusterName := viper.GetString("cluster-name")
		githubHost := viper.GetString("github.host")
		githubOwner := viper.GetString("github.owner")
		kubefirstStateStoreBucket := viper.GetString("kubefirst.state-store.bucket")

		//! computed
		chartmuseumIngressUrl := fmt.Sprintf("https://chartmuseum.%s", awsHostedZoneName)
		chartmuseumIngressUrlNoHttps := fmt.Sprintf("chartmuseum.%s", awsHostedZoneName)
		argocdIngressUrl := fmt.Sprintf("https://argocd.%s", awsHostedZoneName)
		argocdIngressUrlNoHttps := fmt.Sprintf("argocd.%s", awsHostedZoneName)
		argoIngressUrl := fmt.Sprintf("https://argo.%s", awsHostedZoneName)
		argoIngressUrlNoHttps := fmt.Sprintf("argo.%s", awsHostedZoneName)
		gitopsUrlNoHttps := fmt.Sprintf("github.com/%s/gitops.git", viper.GetString("github.owner"))
		vaultIngressUrl := fmt.Sprintf("https://vault.%s", awsHostedZoneName)
		vaultIngressUrlNoHttps := fmt.Sprintf("vault.%s", awsHostedZoneName)
		vouchIngressUrl := fmt.Sprintf("https://vouch.%s", awsHostedZoneName)
		kubefirstIngressUrl := fmt.Sprintf("kubefirst.%s", awsHostedZoneName)
		atlantisIngressUrlNoHttps := fmt.Sprintf("atlantis.%s", awsHostedZoneName)
		atlantisIngressUrl := fmt.Sprintf("https://atlantis.%s", awsHostedZoneName)
		gitlabIngressUrl := fmt.Sprintf("https://gitlab.%s", awsHostedZoneName)

		// todo consolidate
		metaphorDevelopmentIngressUrlNoHttps := fmt.Sprintf("metaphor-development.%s", awsHostedZoneName)
		metaphorStagingIngressUrlNoHttps := fmt.Sprintf("metaphor-staging.%s", awsHostedZoneName)
		metaphorProductionIngressUrlNoHttps := fmt.Sprintf("metaphor-production.%s", awsHostedZoneName)
		metaphorDevelopmentIngressUrl := fmt.Sprintf("https://metaphor-development.%s", awsHostedZoneName)
		metaphorStagingIngressUrl := fmt.Sprintf("https://metaphor-staging.%s", awsHostedZoneName)
		metaphorProductionIngressUrl := fmt.Sprintf("https://metaphor-production.%s", awsHostedZoneName)
		// todo consolidate
		metaphorFrontendDevelopmentIngressUrlNoHttps := fmt.Sprintf("metaphor-frontend-development.%s", awsHostedZoneName)
		metaphorFrontendStagingIngressUrlNoHttps := fmt.Sprintf("metaphor-frontend-staging.%s", awsHostedZoneName)
		metaphorFrontendProductionIngressUrlNoHttps := fmt.Sprintf("metaphor-frontend-production.%s", awsHostedZoneName)
		metaphorFrontendDevelopmentIngressUrl := fmt.Sprintf("https://metaphor-frontend-development.%s", awsHostedZoneName)
		metaphorFrontendStagingIngressUrl := fmt.Sprintf("https://metaphor-frontend-staging.%s", awsHostedZoneName)
		metaphorFrontendProductionIngressUrl := fmt.Sprintf("https://metaphor-frontend-production.%s", awsHostedZoneName)
		// todo consolidate
		metaphorGoDevelopmentIngressUrlNoHttps := fmt.Sprintf("metaphor-go-development.%s", awsHostedZoneName)
		metaphorGoStagingIngressUrlNoHttps := fmt.Sprintf("metaphor-go-staging.%s", awsHostedZoneName)
		metaphorGoProductionIngressUrlNoHttps := fmt.Sprintf("metaphor-go-production.%s", awsHostedZoneName)
		metaphorGoDevelopmentIngressUrl := fmt.Sprintf("https://metaphor-go-development.%s", awsHostedZoneName)
		metaphorGoStagingIngressUrl := fmt.Sprintf("https://metaphor-go-staging.%s", awsHostedZoneName)
		metaphorGoProductionIngressUrl := fmt.Sprintf("https://metaphor-go-production.%s", awsHostedZoneName)

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
		newContents = strings.Replace(newContents, "<GITHUB_OWNER>", githubOwner, -1)
		newContents = strings.Replace(newContents, "<GITHUB_HOST>", githubHost, -1)
		newContents = strings.Replace(newContents, "<ARGOCD_INGRESS_URL_NO_HTTPS>", argocdIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<KUBEFIRST_INGRESS_URL_NO_HTTPS>", kubefirstIngressUrl, -1)
		newContents = strings.Replace(newContents, "<CHARTMUSEUM_INGRESS_URL>", chartmuseumIngressUrl, -1)
		newContents = strings.Replace(newContents, "<CHARTMUSEUM_INGRESS_URL_NO_HTTPS>", chartmuseumIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<CHARTMUSEUM_STORAGE_BUCKET_NAME>", chartmuseumStorageBucketName, -1)
		newContents = strings.Replace(newContents, "<CLUSTER_NAME>", clusterName, -1)
		newContents = strings.Replace(newContents, "<AWS_HOSTED_ZONE_NAME>", awsHostedZoneName, -1)
		newContents = strings.Replace(newContents, "<VAULT_INGRESS_URL>", vaultIngressUrlNoHttps, -1)

		// todo consolidate this?
		newContents = strings.Replace(newContents, "<METAPHOR_DEVELOPMENT_INGRESS_URL_NO_HTTPS>", metaphorDevelopmentIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_STAGING_INGRESS_URL_NO_HTTPS>", metaphorStagingIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_PRODUCTION_INGRESS_URL_NO_HTTPS>", metaphorProductionIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_DEVELOPMENT_INGRESS_URL>", metaphorDevelopmentIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_STAGING_INGRESS_URL>", metaphorStagingIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_PRODUCTION_INGRESS_URL>", metaphorProductionIngressUrl, -1)

		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_DEVELOPMENT_INGRESS_URL_NO_HTTPS>", metaphorFrontendDevelopmentIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_STAGING_INGRESS_URL_NO_HTTPS>", metaphorFrontendStagingIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_PRODUCTION_INGRESS_URL_NO_HTTPS>", metaphorFrontendProductionIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_DEVELOPMENT_INGRESS_URL>", metaphorFrontendDevelopmentIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_STAGING_INGRESS_URL>", metaphorFrontendStagingIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_FRONTEND_PRODUCTION_INGRESS_URL>", metaphorFrontendProductionIngressUrl, -1)

		newContents = strings.Replace(newContents, "<METAPHOR_GO_DEVELOPMENT_INGRESS_URL_NO_HTTPS>", metaphorGoDevelopmentIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_GO_STAGING_INGRESS_URL_NO_HTTPS>", metaphorGoStagingIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_GO_PRODUCTION_INGRESS_URL_NO_HTTPS>", metaphorGoProductionIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_GO_DEVELOPMENT_INGRESS_URL>", metaphorGoDevelopmentIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_GO_STAGING_INGRESS_URL>", metaphorGoStagingIngressUrl, -1)
		newContents = strings.Replace(newContents, "<METAPHOR_GO_PRODUCTION_INGRESS_URL>", metaphorGoProductionIngressUrl, -1)
		newContents = strings.Replace(newContents, "<KUBEFIRST_VERSION>", "TODO", -1) // todo get version

		//! terraform
		// ? argocd ingress url might be in registry?
		newContents = strings.Replace(newContents, "<ARGOCD_INGRESS_URL>", argocdIngressUrl, -1)

		// didnt see
		newContents = strings.Replace(newContents, "<ARGO_WORKFLOWS_INGRESS_URL>", argoIngressUrl, -1)
		newContents = strings.Replace(newContents, "<ARGO_WORKFLOWS_INGRESS_URL_NO_HTTPS>", argoIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<VAULT_INGRESS_URL>", vaultIngressUrl, -1)
		newContents = strings.Replace(newContents, "<VOUCH_INGRESS_URL>", vouchIngressUrl, -1)
		newContents = strings.Replace(newContents, "<ATLANTIS_WEBHOOK_URL>", atlantisWebhookUrl, -1)
		newContents = strings.Replace(newContents, "<ATLANTIS_INGRESS_URL_NO_HTTPS>", atlantisIngressUrlNoHttps, -1)
		newContents = strings.Replace(newContents, "<ATLANTIS_INGRESS_URL>", atlantisIngressUrl, -1)
		newContents = strings.Replace(newContents, "<GITLAB_INGRESS_URL>", gitlabIngressUrl, -1)
		newContents = strings.Replace(newContents, "<KUBEFIRST_STATE_STORE_BUCKET>", kubefirstStateStoreBucket, -1)
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
