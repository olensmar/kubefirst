package aws

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/internal/downloadManager"
	"github.com/kubefirst/kubefirst/internal/gitClient"
	"github.com/kubefirst/kubefirst/internal/helm"
	"github.com/kubefirst/kubefirst/internal/reports"
	"github.com/kubefirst/kubefirst/internal/terraform"
	"github.com/kubefirst/kubefirst/internal/wrappers"
	"github.com/kubefirst/kubefirst/pkg"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runAws(cmd *cobra.Command, args []string) error {

	config := configs.ReadConfig()

	//* confirm with user to continue
	var userInput string
	printConfirmationScreen()
	go counter()
	fmt.Println("to proceed, type 'yes' any other answer will exit")
	fmt.Scanln(&userInput)
	fmt.Println("proceeding with cluster create")
	// if userInput != "yes" {
	// 	os.Exit(1)
	// }
	//* confirm with user to continue

	silentMode := false
	dryRun := false
	awsHostedZone := viper.GetString("aws.hosted-zone-name")
	gitopsTemplateBranch := viper.GetString("template-repo.gitops.branch")
	gitopsTemplateUrl := viper.GetString("template-repo.gitops.url")

	//* emit cluster install started
	if useTelemetryFlag {
		if err := wrappers.SendSegmentIoTelemetry(awsHostedZone, pkg.MetricMgmtClusterInstallStarted); err != nil {
			log.Println(err)
		}
	}

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
	}
	//* git clone and detokenize the gitops repository
	if !viper.GetBool("template-repo.gitops.cloned") {

		//* step 1
		pkg.InformUser("generating your new gitops repository", silentMode)

		gitopsRepoDir := fmt.Sprintf("%s/%s", config.K1FolderPath, "gitops")
		gitClient.CloneRepo(gitopsTemplateUrl, gitopsTemplateBranch, gitopsRepoDir)
		log.Println("gitops repository creation complete")

		//* step 2
		// adjust content in gitops repository
		opt := cp.Options{
			Skip: func(src string) (bool, error) {
				if strings.HasSuffix(src, ".git") {
					return true, nil
				} else if strings.Index(src, "/.terraform") > 0 {
					return true, nil
				}
				//Add more stuff to be ignored here
				return false, nil

			},
		}

		// clear out the root of `gitops-template` once we move
		// all the content we only remove the different root folders
		os.RemoveAll(gitopsRepoDir + "/components")
		os.RemoveAll(gitopsRepoDir + "/localhost")
		os.RemoveAll(gitopsRepoDir + "/registry")
		os.RemoveAll(gitopsRepoDir + "/validation")
		os.RemoveAll(gitopsRepoDir + "/terraform")
		os.RemoveAll(gitopsRepoDir + "/.gitignore")
		os.RemoveAll(gitopsRepoDir + "/LICENSE")
		os.RemoveAll(gitopsRepoDir + "/README.md")
		os.RemoveAll(gitopsRepoDir + "/atlantis.yaml")
		os.RemoveAll(gitopsRepoDir + "/logo.png")

		driverContent := fmt.Sprintf("%s/%s-%s", gitopsRepoDir, viper.GetString("cloud-provider"), viper.GetString("git-provider"))
		err := cp.Copy(driverContent, gitopsRepoDir, opt)
		if err != nil {
			log.Println("Error populating gitops with local setup:", err)
			return err
		}
		os.RemoveAll(driverContent)

		//* step 3 -- gitClient.CommitAndPush -- warning origin is github
		pkg.DetokenizeV2(gitopsRepoDir)

		//* step 4 add a new remote of the github user who's token we have
		repo, err := git.PlainOpen(gitopsRepoDir)
		if err != nil {
			log.Print("error opening repo at:", gitopsRepoDir)
		}
		destinationGitopsRepoURL := viper.GetString("github.repo.gitops.giturl")
		log.Printf("git remote add github %s", destinationGitopsRepoURL)
		_, err = repo.CreateRemote(&gitConfig.RemoteConfig{
			Name: "github",
			URLs: []string{destinationGitopsRepoURL},
		})
		if err != nil {
			log.Panicf("Error creating remote %s at: %s - %s", viper.GetString("git-provider"), destinationGitopsRepoURL, err)
		}

		//* step 5 commit newly detokenized content
		w, _ := repo.Worktree()

		log.Printf("committing detokenized %s content", "gitops")
		status, err := w.Status()
		if err != nil {
			log.Println("error getting worktree status", err)
		}

		for file, _ := range status {
			_, err = w.Add(file)
			if err != nil {
				log.Println("error getting worktree status", err)
			}
		}
		w.Commit(fmt.Sprintf("[ci skip] committing initial detokenized %s content", destinationGitopsRepoURL), &git.CommitOptions{
			Author: &object.Signature{
				Name:  "kubefirst-bot",
				Email: "kubefirst-bot@kubefirst.com",
				When:  time.Now(),
			},
		}) // todo emit init telemetry end

		viper.Set("template-repo.gitops.cloned", true)
		viper.WriteConfig()
	} else {
		log.Println("gitops repository generation already complete - continuing")
	}

	//* terraform apply github repositories
	executionControl := viper.GetBool("terraform.github.apply.complete")
	// create github teams in the org and repositories
	if !executionControl {
		pkg.InformUser("creating github resources with terraform", silentMode)

		tfEntrypoint := config.TerraformGithubEntrypointPath
		err := terraform.InitApplyAutoApprove(dryRun, tfEntrypoint)
		if err != nil {
			log.Printf("error executing terraform apply %s", tfEntrypoint)
			return err
		}

		viper.Set("terraform.github.apply.complete", true)
		viper.WriteConfig()

		pkg.InformUser(fmt.Sprintf("created github repositories in github.com/%s", viper.GetString("github.owner")), silentMode)
		log.Println(fmt.Sprintf("created github repositories in https://github.com/%s", viper.GetString("github.owner")))
		log.Printf(fmt.Sprintf("     %s", viper.GetString("github.repo.gitops.url")))
		log.Printf(fmt.Sprintf("     %s", viper.GetString("github.repo.metaphor.url")))
		log.Printf(fmt.Sprintf("     %s", viper.GetString("github.repo.metaphor-frontend.url")))
		log.Printf(fmt.Sprintf("     %s", viper.GetString("github.repo.metaphor-go.url")))
		// progressPrinter.IncrementTracker("step-github", 1)
	} else {
		log.Println("already created github terraform resources")
	}

	//* terraform apply aws cloud resources
	executionControl = viper.GetBool("terraform.aws.apply.complete")
	// create github teams in the org and repositories
	if !executionControl {
		pkg.InformUser("Creating aws resources with terraform", silentMode)

		tfEntrypoint := config.TerraformAwsEntrypointPath
		err := terraform.InitApplyAutoApprove(dryRun, tfEntrypoint)
		if err != nil {
			log.Printf("error executing terraform apply %s", tfEntrypoint)
			return err
		}

		viper.Set("terraform.aws.apply.complete", true)
		viper.WriteConfig()

		pkg.InformUser(fmt.Sprintf("Created aws cloud resources for kubefirst cluster"), silentMode)
		// progressPrinter.IncrementTracker("step-github", 1)
	} else {
		log.Println("already created aws terraform resources")
	}

	//* terraform output kms key from terraform/aws and detokenize, commit, push
	executionControl = viper.GetBool("template-repo.gitops.pushed")
	// create github teams in the org and repositories
	if !executionControl {
		pkg.InformUser("adding kms id to gitops repository", silentMode)

		//* get kms key id terraform.OutputSingleValue
		tfEntrypoint := config.TerraformAwsEntrypointPath
		tfOutputName := "vault_unseal_kms_key"
		kmsKeyId, err := terraform.OutputSingleValue(dryRun, tfEntrypoint, tfOutputName)
		if err != nil {
			log.Printf("error getting terraform output %s value %s", tfOutputName, kmsKeyId)
			return err
		}
		log.Printf("terraform output: %s = %s", tfOutputName, kmsKeyId)

		viper.Set("terraform.aws.outputs.kms-key.id", kmsKeyId)
		viper.WriteConfig()

		//* step 3 -- gitClient.CommitAndPush -- warning origin is github
		//* detokenize the kms key id for vault and dynamodb table name
		pkg.DetokenizeV2(config.GitOpsRepoPath)

		//* step 4 add a new remote of the github user who's token we have
		repo, err := git.PlainOpen(config.GitOpsRepoPath)
		if err != nil {
			log.Print("error opening repo at:", config.GitOpsRepoPath)
		}
		//* step 5 commit newly detokenized content
		w, _ := repo.Worktree()

		log.Printf("committing detokenized %s content", "gitops")
		status, err := w.Status()
		if err != nil {
			log.Println("error getting worktree status", err)
		}

		for file, _ := range status {
			_, err = w.Add(file)
			if err != nil {
				log.Println("error getting worktree status", err)
			}
		}
		w.Commit("[ci skip] committing detokenized kms key id", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "kubefirst-bot",
				Email: "kubefirst-bot@kubefirst.com",
				When:  time.Now(),
			},
		})

		token := os.Getenv("GITHUB_TOKEN")
		if len(token) == 0 {
			log.Println("no GITHUB_TOKEN provided, unable to push to github")
		}

		err = repo.Push(&git.PushOptions{
			RemoteName: "github", // todo parameter
			Auth: &http.BasicAuth{
				Username: "kubefirst-bot",
				Password: token,
			},
		})
		if err != nil {
			log.Panicf("error pushing to remote github: %s", err)
		}
		log.Println("successfully pushed detokenized gitops content to github/", viper.GetString("github.owner"))
		viper.Set("github.gitops.hydrated", true)
		viper.WriteConfig()

		viper.Set("template-repo.gitops.pushed", true)
		viper.WriteConfig()

		pkg.InformUser(fmt.Sprintf("Created github repos in github.com/%s", viper.GetString("github.owner")), silentMode)
		// progressPrinter.IncrementTracker("step-github", 1)
	} else {
		log.Println("already created github terraform resources")
	}

	// todo restore ssl... also automatically backup ssl at the end
	// informUser("Attempt to recycle certs", globalFlags.SilentMode)
	// restoreSSLCmd.RunE(cmd, args)

	// todo create initial argocd repository (this is the connection to argocd as a 'repo') destinationGitopsRepoUrl
	//* investigate - is this doing the same thing as pkg.CreateSSHKey where it writes a file?
	// gitopsRepo := fmt.Sprintf("git@github.com:%s/gitops.git", viper.GetString("github.owner"))
	// argocd.CreateInitialArgoCDRepository(gitopsRepo)

	// helm add argo repository && update
	// todo move to config.?
	helmRepo := helm.HelmRepo{
		RepoName:     pkg.HelmRepoName,
		RepoURL:      pkg.HelmRepoURL,
		ChartName:    pkg.HelmRepoChartName,
		Namespace:    pkg.HelmRepoNamespace,
		ChartVersion: pkg.HelmRepoChartVersion,
	}
	//* helm repo add argocd
	executionControl = viper.GetBool("argocd.helm.repo.updated")
	if !executionControl {
		pkg.InformUser(fmt.Sprintf("helm repo add %s %s and helm repo update", helmRepo.RepoName, helmRepo.RepoURL), silentMode)
		helm.AddRepoAndUpdateRepo(dryRun, helmRepo)
	}

	//* helm install argocd
	executionControl = viper.GetBool("argocd.helm.install.complete")
	if !executionControl {
		pkg.InformUser(fmt.Sprintf("helm install %s and wait", helmRepo.RepoName), silentMode)
		helm.Install(dryRun, helmRepo)
	}
	log.Println("hell yes")

	// todo wait for argocd to be ready

	// todo set argocd credentials

	//! stop here before continuing
	// todo apply argocd registry

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
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.gitops.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.gitops.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor-frontend.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor-frontend.branch")))
	createKubefirstSummary.WriteString(fmt.Sprintf("  %s\n", viper.GetString("template-repo.metaphor-go.url")))
	createKubefirstSummary.WriteString(fmt.Sprintf("    branch:  %s\n", viper.GetString("template-repo.metaphor-go.branch")))

	fmt.Println(reports.StyleMessage(createKubefirstSummary.String()))
}
