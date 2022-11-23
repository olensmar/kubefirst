package local

import (
	"github.com/spf13/cobra"
)

var (
	// enableConsole bool
	gitProvider string
)

func NewCommand() *cobra.Command {

	awsCmd := &cobra.Command{
		Use:     "aws",
		Short:   "kubefirst aws installation",
		Long:    "kubefirst aws",
		PreRunE: validateAws,
		RunE:    runAws,
		// PostRunE: runPostAws,
	}

	// awsCmd.Flags().BoolVar(&enableConsole, "enable-console", true, "If hand-off screen will be presented on a browser UI")
	awsCmd.Flags().StringVar(&gitProvider, "git-provider", "github", "the git provider to use. (i.e. gitlab|github)")

	// on error, doesnt show helper/usage
	awsCmd.SilenceUsage = true

	// wire up new commands
	// awsCmd.AddCommand(NewCommandConnect())

	return awsCmd
}

func runAws(cmd *cobra.Command, args []string) error {

	// config := configs.ReadConfig()
	// gitProviderFlag, err := cmd.Flags().GetString("git-provider")
	// if err != nil {
	// 	return err
	// }
	// log.Println("git-provider flag value", gitProviderFlag)

	return nil

}
