package aws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cip8/autoname"
	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/aws"
	"github.com/kubefirst/kubefirst/internal/handlers"
	"github.com/kubefirst/kubefirst/internal/services"
	"github.com/kubefirst/kubefirst/internal/wrappers"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateAws is responsible for gathering all of the information required to execute a kubefirst aws cloud creation with github (currently)
// this function needs to provide all the generated values and provides a single space for writing and updating configuration up front.
func validateAws(cmd *cobra.Command, args []string) error {

	// infoCmd.Run(cmd, args)

	config := configs.ReadConfig()

	adminEmailFlag, err := cmd.Flags().GetString("admin-email")
	if err != nil {
		log.Panic(err)
	}

	awsProfileFlag, err := cmd.Flags().GetString("aws-profile")
	if err != nil {
		log.Panic(err)
	}

	awsRegionFlag, err := cmd.Flags().GetString("aws-region")
	if err != nil {
		log.Panic(err)
	}
	if awsRegionFlag == "" {
		//* if the region is not set we want to force the sdk to look at
		//* $HOME/.aws/config and use the region set in the users config
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
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
	//! hack
	// gitopsTemplateBranchFlag, err := cmd.Flags().GetString("gitops-template-branch")
	// if err != nil {
	// 	log.Panic(err)
	// }

	gitProviderFlag, err := cmd.Flags().GetString("git-provider")
	if err != nil {
		log.Panic(err)
	}

	awsHostedZoneNameFlag, err := cmd.Flags().GetString("aws-hosted-zone-name")
	if err != nil {
		log.Panic(err)
	}

	kbotPassword, err := cmd.Flags().GetString("kbot-password")
	if err != nil {
		log.Panic(err)
	}

	if strings.HasSuffix(awsHostedZoneNameFlag, ".") {
		awsHostedZoneNameFlag = awsHostedZoneNameFlag[:len(awsHostedZoneNameFlag)-1]
	}

	useTelemetryFlag, err := cmd.Flags().GetBool("use-telemetry")
	if err != nil {
		log.Panic(err)
	}

	if useTelemetryFlag {
		if err := wrappers.SendSegmentIoTelemetry(awsHostedZoneNameFlag, pkg.MetricInitStarted); err != nil {
			log.Println(err)
		}
	}
	//! hack
	// if err := pkg.ValidateK1Folder(config.K1FolderPath); err != nil {
	// 	return err
	// }

	// todo validate flags
	viper.Set("admin-email", adminEmailFlag)
	viper.Set("argocd.local.service", config.LocalArgoCdURL)
	viper.Set("cloud-provider", cloudProviderFlag)
	viper.Set("git-provider", gitProviderFlag)

	//! hack
	// viper.Set("gitops-template.repo.branch", gitopsTemplateBranchFlag)
	viper.Set("template-repo.gitops.branch", "domain-refactor")
	viper.Set("template-repo.gitops.url", gitopsTemplateUrlFlag)
	// todo accommodate metaphor branch and repo override more intelligently
	viper.Set("template-repo.metaphor.url", fmt.Sprintf("https://github.com/%s/metaphor.git", "kubefirst"))
	viper.Set("template-repo.metaphor.branch", "main")
	viper.Set("template-repo.metaphor-frontend.url", fmt.Sprintf("https://github.com/%s/metaphor-frontend.git", "kubefirst"))
	viper.Set("template-repo.metaphor-frontend.branch", "main")
	viper.Set("template-repo.metaphor-go.url", fmt.Sprintf("https://github.com/%s/metaphor-go.git", "kubefirst"))
	viper.Set("template-repo.metaphor-go.branch", "main")
	viper.Set("github.atlantis.webhook.secret", pkg.Random(20))
	viper.Set("github.atlantis.webhook.url", fmt.Sprintf("https://atlantis.%s/events", awsHostedZoneNameFlag))
	viper.Set("github.repo.gitops.url", fmt.Sprintf("https://github.com/%s/gitops.git", githubOwnerFlag))
	viper.Set("github.repo.metaphor.url", fmt.Sprintf("https://github.com/%s/metaphor.git", githubOwnerFlag))
	viper.Set("github.repo.metaphor-frontend.url", fmt.Sprintf("https://github.com/%s/metaphor-frontend.git", githubOwnerFlag))
	viper.Set("github.repo.metaphor-go.url", fmt.Sprintf("https://github.com/%s/metaphor-go.git", githubOwnerFlag))
	viper.WriteConfig()

	//* github checks
	executionControl := viper.GetBool("kubefirst.checks.github.complete")
	if !executionControl {
		httpClient := http.DefaultClient
		githubToken := config.GithubToken
		if len(githubToken) == 0 {
			return errors.New("ephemeral tokens not supported for cloud installations, please set a GITHUB_TOKEN environment variable to continue\n https://docs.kubefirst.io/kubefirst/github/install.html#step-3-kubefirst-init")
		}
		gitHubService := services.NewGitHubService(httpClient)
		gitHubHandler := handlers.NewGitHubHandler(gitHubService)

		// get Github data to set user based on the provided token
		log.Println("verifying github user")
		githubUser, err := gitHubHandler.GetGitHubUser(githubToken)
		if err != nil {
			return err
		}
		log.Println("github user is: ", githubUser)
		// todo evaluate if cloudProviderFlag == "local" {githubOwnerFlag = githubUser} and the rest of the execution is the same

		err = gitHubHandler.CheckGithubOrganizationPermissions(githubToken, githubOwnerFlag, githubUser)
		if err != nil {
			return err
		}
		viper.Set("github.owner", githubOwnerFlag)
		viper.Set("github.user", githubUser)
		viper.Set("kubefirst.checks.github.complete", true)
		viper.WriteConfig()
	} else {
		log.Println("already completed github checks - continuing")
	}

	executionControl = viper.GetBool("kubefirst.checks.aws.complete")
	if !executionControl {
		log.Println("getting aws account information")
		awsAccountId, awsIamArn, awsRegion, err := aws.GetAccountInfoV2(awsProfileFlag, awsRegionFlag)
		if err != nil {
			return err
		}
		log.Printf("aws account id: %s\naws user arn: %s", awsAccountId, awsIamArn)

		log.Println("getting aws hosted zone id for zone ", awsHostedZoneNameFlag)
		awsHostedZoneId := aws.GetHostedZoneId(awsProfileFlag, awsRegion, awsHostedZoneNameFlag)
		log.Printf("aws hosted zone id %s", awsHostedZoneId)

		log.Printf("creating state store bucket ")
		randomName := strings.ReplaceAll(autoname.Generate(), "_", "-")

		kubefirstStateStoreBucketName := fmt.Sprintf("k1-state-store-%s", randomName)
		err = aws.CreateS3Bucket(awsProfileFlag, awsRegion, clusterNameFlag, kubefirstStateStoreBucketName)
		if err != nil {
			log.Printf("creating state store bucket ")
			return err
		}
		viper.Set("kubefirst.state-store.bucket", kubefirstStateStoreBucketName)
		viper.Set("kubefirst.bucket.random-name", randomName)
		viper.Set("kubefirst.telemetry", useTelemetryFlag)
		viper.Set("cluster-name", clusterNameFlag)
		viper.Set("vault.local.service", config.LocalVaultURL)
		viper.Set("aws.account-id", awsAccountId)
		viper.Set("aws.iam-arn", awsIamArn)
		viper.Set("aws.hosted-zone-id", awsHostedZoneId)
		viper.Set("aws.hosted-zone-name", awsHostedZoneNameFlag)
		viper.Set("aws.profile", awsProfileFlag)
		viper.Set("aws.region", awsRegion)
		viper.Set("kubefirst.checks.aws.complete", true)
		viper.WriteConfig()
	} else {
		log.Println("already completed aws checks - continuing")
	}

	executionControl = viper.GetBool("kubefirst.checks.kbot.complete")
	if !executionControl {
		log.Println("creating an ssh key pair for your new cloud infrastructure")
		sshPrivateKey, sshPublicKey, err := pkg.CreateSshKeyPair()
		if err != nil {
			return err
		}
		if len(kbotPassword) == 0 {
			kbotPassword = pkg.Random(20)
		}
		log.Println("ssh key pair creation complete")
		githubOwnerRootGitUrl := fmt.Sprintf("git@github.com:%s", githubOwnerFlag)
		viper.Set("kubefirst.bot.password", kbotPassword)
		viper.Set("kubefirst.bot.private-key", sshPrivateKey)
		viper.Set("kubefirst.bot.public-key", sshPublicKey)
		viper.Set("kubefirst.bot.user", "kbot")
		viper.Set("github.repo.gitops.giturl", fmt.Sprintf("%s/gitops.git", githubOwnerRootGitUrl))
		viper.Set("kubefirst.checks.kbot.complete", true)
		viper.WriteConfig()
		// todo, is this a hangover from initial gitlab? do we need this?
		log.Println("creating argocd-init-values.yaml for initial install")
		//* ex: `git@github.com:kubefirst` this is allows argocd access to the kubefirst organization repos
		err = pkg.WriteGithubArgoCdInitValuesFile(githubOwnerRootGitUrl, sshPrivateKey)
		if err != nil {
			return err
		}
		log.Println("argocd-init-values.yaml creation complete")
	}

	log.Println("validation and kubefirst cli environment check is complete")

	if useTelemetryFlag {
		if err := wrappers.SendSegmentIoTelemetry(awsHostedZoneNameFlag, pkg.MetricInitCompleted); err != nil {
			log.Println(err)
		}
	}

	// todo progress bars
	// time.Sleep(time.Millisecond * 100) // to allow progress bars to finish

	return nil
}
