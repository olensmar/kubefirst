package create

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
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

func runCreate(cmd *cobra.Command, args []string) error {

	fmt.Println("run new create Cmd")

	return errors.New("fake error")
}
