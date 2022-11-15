package create

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	useTelemetry   bool
	dryRun         bool
	silentMode     bool
	enableConsole  bool
	gitOpsBranch   string
	gitOpsRepo     string
	awsHostedZone  string
	metaphorBranch string
	adminEmail     string
	templateTag    string
)

func NewCommand() *cobra.Command {

	createCmd := &cobra.Command{
		Use:   "civo", // rebrand this to cluster create?
		Short: "Kubefirst localhost installation",
		Long:  "Kubefirst localhost enable a localhost installation without the requirement of a cloud provider.",
		// PreRunE:  validateCreate,
		RunE: runCreate,
		// PostRunE: runCreate,
	}

	// createCmd.Flags().BoolVar(&useTelemetry, "use-telemetry", true, "installer will not send telemetry about this installation")
	// createCmd.Flags().StringVar(&adminEmail, "admin-email", "", "the email address for the administrator as well as for lets-encrypt certificate emails")

	// on error, doesnt show helper/usage
	createCmd.SilenceUsage = true

	return createCmd
}

type CivoExecutionControl struct {
	Step1 string
	Step2 string
}

func runCreate(cmd *cobra.Command, args []string) error {

	executionControl := CivoExecutionControl{}
	executionControl.Step1 = "kubefirst.init.env-check.complete"

	//* 1. civo token - validate and require
	//* 2. github token - validate and require
	//*
	//*

	silentMode := false

	if !viper.GetBool(executionControl.Step1) {
		pkg.InformUser("validating `kubefirst init` environment", silentMode)
		log.Println("INIT: checking environment variables")

		civoToken := os.Getenv("CIVO_TOKEN") //* todo fix GITHUB_AUTH_TOKEN
		if civoToken == "" {
			log.Println("Unauthorized: No CIVO_TOKEN environment variable is present.")
			return fmt.Errorf("unauthorized: missing CIVO_TOKEN environment variable")
		}
		log.Println("INIT: CIVO_TOKEN is set")

		// todo add the ephemeral token logic here
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			log.Println("Unauthorized: No GITHUB_TOKEN environment variable is present.")
			return fmt.Errorf("unauthorized: missing GITHUB_TOKEN environment variable")
		}
		log.Println("INIT: GITHUB_TOKEN is set")

	} else {
		log.Println("INIT: environment variables all set - continuing")
		log.Println("executionControl.Step1 complete")
	}

	return errors.New("\n\nno error: get the next rock")
}
