package aws

import "github.com/spf13/cobra"

var (
	adminEmailFlag           string
	cloudProviderFlag        string
	clusterNameFlag          string
	githubOwner              string
	gitopsTemplateUrlFlag    string
	gitopsTemplateBranchFlag string
	gitProviderFlag          string
	hostedZoneNameFlag       string
	useTelemetryFlag         bool
)

func NewCommand() *cobra.Command {

	awsCmd := &cobra.Command{
		Use:     "aws",
		Short:   "kubefirst aws installation",
		Long:    "kubefirst aws",
		PreRunE: validateAws, // todo what should this function be called?
		RunE:    runAws,
		// PostRunE: runPostAws,
	}

	// todo review defaults and update descriptions
	awsCmd.Flags().StringVar(&adminEmailFlag, "admin-email", "jared@kubeshop.io", "email address for let's encrypt certificate notifications")
	awsCmd.Flags().StringVar(&cloudProviderFlag, "cloud-provider", "aws", "the git provider to use. (i.e. gitlab|github)")
	awsCmd.Flags().StringVar(&clusterNameFlag, "cluster-name", "kubefirst", "the name of the cluster to create")
	awsCmd.Flags().StringVar(&githubOwner, "github-owner", "your-dns-io", "the GitHub owner of the new gitops and metaphor repositories")
	// awsCmd.MarkFlagRequired("github-owner")
	awsCmd.Flags().StringVar(&gitopsTemplateBranchFlag, "gitops-template-branch", "main", "the branch to clone for the gitops-template repository")
	awsCmd.Flags().StringVar(&gitopsTemplateUrlFlag, "gitops-template-url", "https://github.com/kubefirst/gitops-template.git", "the fully qualified url to the gitops-template repository to clone")
	awsCmd.Flags().StringVar(&gitProviderFlag, "git-provider", "github", "the git provider to use. (i.e. gitlab|github)")
	awsCmd.Flags().StringVar(&hostedZoneNameFlag, "hosted-zone-name", "kubernickels.com", "the AWS Hosted Zone to use for DNS records (i.e. your-domain.com|subdomain.your-domain.com)")
	awsCmd.Flags().BoolVar(&useTelemetryFlag, "use-telemetry", true, "whether to emit telemetry")

	// on error, doesnt show helper/usage
	awsCmd.SilenceUsage = true

	// wire up new commands
	// awsCmd.AddCommand(NewCommandConnect())

	return awsCmd
}
