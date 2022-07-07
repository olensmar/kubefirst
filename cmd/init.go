package cmd

import (
	"fmt"
	"github.com/kubefirst/nebulous/configs"
	"github.com/kubefirst/nebulous/internal/aws"
	"github.com/kubefirst/nebulous/internal/downloadManager"
	"github.com/kubefirst/nebulous/internal/gitlab"
	"github.com/kubefirst/nebulous/internal/telemetry"
	"github.com/kubefirst/nebulous/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		config := configs.ReadConfig()

		var err error
		config.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			panic(err)
		}

		log.Println("dry run enabled:", config.DryRun)

		pkg.SetupProgress(10)
		trackers := pkg.GetTrackers()
		trackers[pkg.TrackerStage0] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage0, 1)}
		trackers[pkg.TrackerStage1] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage1, 1)}
		trackers[pkg.TrackerStage2] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage2, 1)}
		trackers[pkg.TrackerStage3] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage3, 1)}
		trackers[pkg.TrackerStage4] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage4, 1)}
		trackers[pkg.TrackerStage5] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage5, 3)}
		trackers[pkg.TrackerStage6] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage6, 1)}
		trackers[pkg.TrackerStage7] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage7, 4)}
		trackers[pkg.TrackerStage8] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage8, 1)}
		trackers[pkg.TrackerStage9] = &pkg.ActionTracker{Tracker: pkg.CreateTracker(pkg.TrackerStage9, 1)}
		infoCmd.Run(cmd, args)
		hostedZoneName, _ := cmd.Flags().GetString("hosted-zone-name")
		metricName := "kubefirst.init.started"
		metricDomain := hostedZoneName

		if !config.DryRun {
			telemetry.SendTelemetry(metricDomain, metricName)
		} else {
			log.Printf("[#99] Dry-run mode, telemetry skipped:  %s", metricName)
		}

		// todo need to check flags and create config

		// hosted zone name:
		// name of the hosted zone to be used for the kubefirst install
		// if suffixed with a dot (eg. kubefirst.com.), the dot will be stripped
		if strings.HasSuffix(hostedZoneName, ".") {
			hostedZoneName = hostedZoneName[:len(hostedZoneName)-1]
		}
		log.Println("hostedZoneName:", hostedZoneName)
		viper.Set("aws.hostedzonename", hostedZoneName)
		viper.WriteConfig()
		// admin email
		// used for letsencrypt notifications and the gitlab root account
		adminEmail, _ := cmd.Flags().GetString("admin-email")
		log.Println("adminEmail:", adminEmail)
		viper.Set("adminemail", adminEmail)

		// region
		// name of the cloud region to provision resources when resources are region-specific
		region, _ := cmd.Flags().GetString("region")
		viper.Set("aws.region", region)
		log.Println("region:", region)

		// hosted zone id
		// so we don't have to keep looking it up from the domain name to use it
		hostedZoneId := aws.GetDNSInfo(hostedZoneName)
		// viper values set in above function
		log.Println("hostedZoneId:", hostedZoneId)
		trackers[pkg.TrackerStage0].Tracker.Increment(1)
		trackers[pkg.TrackerStage1].Tracker.Increment(1)

		// todo: this doesn't default to testing the dns check
		skipHostedZoneCheck := viper.GetBool("init.hostedzonecheck.enabled")
		if !skipHostedZoneCheck {
			log.Println("skipping hosted zone check")
		} else {
			aws.TestHostedZoneLiveness(hostedZoneName, hostedZoneId)
		}
		trackers[pkg.TrackerStage2].Tracker.Increment(1)

		log.Println("calling createSshKeyPair() ")
		createSshKeyPair()
		log.Println("createSshKeyPair() complete")
		trackers[pkg.TrackerStage3].Tracker.Increment(1)

		log.Println("calling cloneGitOpsRepo()")
		cloneGitOpsRepo()
		log.Println("cloneGitOpsRepo() complete")
		trackers[pkg.TrackerStage4].Tracker.Increment(1)

		log.Println("calling download()")
		err = downloadManager.DownloadTools(config, trackers)
		if err != nil {
			panic(err)
		}

		log.Println("download() complete")

		log.Println("calling GetAccountInfo()")
		aws.GetAccountInfo()
		log.Println("GetAccountInfo() complete")
		trackers[pkg.TrackerStage6].Tracker.Increment(1)

		log.Println("calling BucketRand()")
		aws.BucketRand()
		log.Println("BucketRand() complete")

		log.Println("calling detokenize()")
		detokenize(fmt.Sprintf("%s/.kubefirst/gitops", config.HomePath))
		log.Println("detokenize() complete")
		trackers[pkg.TrackerStage8].Tracker.Increment(1)

		// modConfigYaml()
		metricName = "kubefirst.init.completed"

		if !config.DryRun {
			telemetry.SendTelemetry(metricDomain, metricName)
		} else {
			log.Printf("[#99] Dry-run mode, telemetry skipped:  %s", metricName)
		}

		viper.WriteConfig()
		trackers[pkg.TrackerStage9].Tracker.Increment(1)
		time.Sleep(time.Millisecond * 100)
	},
}

func init() {
	config := configs.ReadConfig()
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().String("hosted-zone-name", "", "the domain to provision the kubefirst platform in")
	err := initCmd.MarkFlagRequired("hosted-zone-name")
	if err != nil {
		panic(err)
	}
	initCmd.Flags().String("admin-email", "", "the email address for the administrator as well as for lets-encrypt certificate emails")
	err = initCmd.MarkFlagRequired("admin-email")
	if err != nil {
		panic(err)
	}
	initCmd.Flags().String("cloud", "", "the cloud to provision infrastructure in")
	err = initCmd.MarkFlagRequired("cloud")
	if err != nil {
		panic(err)
	}
	initCmd.Flags().String("region", "", "the region to provision the cloud resources in")
	err = initCmd.MarkFlagRequired("region")
	if err != nil {
		panic(err)
	}
	initCmd.Flags().Bool("clean", false, "delete any local kubefirst content ~/.kubefirst, ~/.flare")

	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	initCmd.PersistentFlags().BoolVarP(&config.DryRun, "dry-run", "s", false, "set to dry-run mode, no changes done on cloud provider selected")
	log.Println("init started")
}

func createSshKeyPair() {
	config := configs.ReadConfig()
	publicKey := viper.GetString("botpublickey")
	if publicKey == "" {
		log.Println("generating new key pair")
		publicKey, privateKey, _ := gitlab.GenerateKey()
		viper.Set("botPublicKey", publicKey)
		viper.Set("botPrivateKey", privateKey)
		err := viper.WriteConfig()
		if err != nil {
			log.Panicf("error: could not write to viper config")
		}
	}
	publicKey = viper.GetString("botpublickey")
	privateKey := viper.GetString("botprivatekey")

	var argocdInitValuesYaml = []byte(fmt.Sprintf(`
server:
  additionalApplications:
  - name: registry
    namespace: argocd
    additionalLabels: {}
    additionalAnnotations: {}
    finalizers:
    - resources-finalizer.argocd.argoproj.io
    project: default
    source:
      repoURL: ssh://soft-serve.soft-serve.svc.cluster.local:22/gitops
      targetRevision: HEAD
      path: registry
    destination:
      server: https://kubernetes.default.svc
      namespace: argocd
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
      - CreateNamespace=true
configs:
  repositories:
    soft-serve-gitops:
      url: ssh://soft-serve.soft-serve.svc.cluster.local:22/gitops
      insecure: 'true'
      type: git
      name: soft-serve-gitops
  credentialTemplates:
    ssh-creds:
      url: ssh://soft-serve.soft-serve.svc.cluster.local:22
      sshPrivateKey: |
        %s
`, strings.ReplaceAll(privateKey, "\n", "\n        ")))

	err := ioutil.WriteFile(fmt.Sprintf("%s/.kubefirst/argocd-init-values.yaml", config.HomePath), argocdInitValuesYaml, 0644)
	if err != nil {
		log.Panicf("error: could not write argocd-init-values.yaml %s", err)
	}
}
