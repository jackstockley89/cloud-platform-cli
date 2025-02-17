package commands

import (
	"errors"
	"os"
	"path/filepath"

	environment "github.com/ministryofjustice/cloud-platform-cli/pkg/environment"
	"github.com/ministryofjustice/cloud-platform-cli/pkg/github"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// variables specific to commands package used to store the values of flags of various environment sub commands
var module, moduleVersion string
var optFlags environment.Options

// skipEnvCheck is a flag to skip the environments repository check.
// This is useful for testing.
var skipEnvCheck bool

// answersFile is a flag to specify the path to the answers file.
var answersFile string

func addEnvironmentCmd(topLevel *cobra.Command) {
	topLevel.AddCommand(environmentCmd)
	environmentCmd.AddCommand(environmentEcrCmd)
	environmentCmd.AddCommand(environmentRdsCmd)
	environmentCmd.AddCommand(environmentS3Cmd)
	environmentCmd.AddCommand(environmentSvcCmd)
	environmentCmd.AddCommand(environmentCreateCmd)
	environmentCmd.AddCommand(environmentPlanCmd)
	environmentCmd.AddCommand(environmentApplyCmd)
	environmentEcrCmd.AddCommand(environmentEcrCreateCmd)
	environmentRdsCmd.AddCommand(environmentRdsCreateCmd)
	environmentS3Cmd.AddCommand(environmentS3CreateCmd)
	environmentSvcCmd.AddCommand(environmentSvcCreateCmd)
	environmentCmd.AddCommand(environmentPrototypeCmd)
	environmentPrototypeCmd.AddCommand(environmentPrototypeCreateCmd)
	environmentCmd.AddCommand(environmentBumpModuleCmd)

	// flags
	environmentApplyCmd.Flags().BoolVar(&optFlags.AllNamespaces, "all-namespaces", false, "Apply all namespaces with -all-namespaces")
	environmentApplyCmd.Flags().StringVarP(&optFlags.Namespace, "namespace", "n", "", "Namespace which you want to perform the apply")
	environmentApplyCmd.Flags().StringVar(&optFlags.CommitTimestamp, "commit-timestamp", "", "Timestamp of current commit from the environment repo")
	// Re-use the environmental variable TF_VAR_github_token to call Github Client which is needed to perform terraform operations on each namespace
	environmentApplyCmd.Flags().StringVar(&optFlags.GithubToken, "github-token", os.Getenv("TF_VAR_github_token"), "Personal access Token from Github ")
	environmentApplyCmd.Flags().StringVar(&optFlags.KubecfgPath, "kubecfg", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path to kubeconfig file")
	environmentApplyCmd.Flags().StringVar(&optFlags.ClusterCtx, "cluster", "", "folder name under namespaces/ inside cloud-platform-environments repo refering to full cluster name")

	// e.g. if this is the Pull rquest to perform the apply: https://github.com/ministryofjustice/cloud-platform-environments/pull/8370, the pr ID is 8370.
	environmentPlanCmd.Flags().IntVar(&optFlags.PRNumber, "prNumber", 0, "Pull request ID or number to which you want to perform the plan")
	environmentPlanCmd.Flags().StringVarP(&optFlags.Namespace, "namespace", "n", "", "Namespace which you want to perform the plan")

	// Re-use the environmental variable TF_VAR_github_token to call Github Client which is needed to perform terraform operations on each namespace
	environmentPlanCmd.Flags().StringVar(&optFlags.GithubToken, "github-token", os.Getenv("TF_VAR_github_token"), "Personal access Token from Github ")
	environmentPlanCmd.Flags().StringVar(&optFlags.KubecfgPath, "kubecfg", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path to kubeconfig file")
	environmentPlanCmd.Flags().StringVar(&optFlags.ClusterCtx, "cluster", "", "folder name under namespaces/ inside cloud-platform-environments repo refering to full cluster name")

	environmentBumpModuleCmd.Flags().StringVarP(&module, "module", "m", "", "Module to upgrade the version")
	environmentBumpModuleCmd.Flags().StringVarP(&moduleVersion, "module-version", "v", "", "Semantic version to bump a module to")

	environmentCreateCmd.Flags().BoolVarP(&skipEnvCheck, "skip-env-check", "s", false, "Skip the environment check")
	environmentCreateCmd.Flags().StringVarP(&answersFile, "answers-file", "a", "", "Path to the answers file")
}

var environmentCmd = &cobra.Command{
	Use:    "environment",
	Short:  `Cloud Platform Environment actions`,
	PreRun: upgradeIfNotLatest,
}

var environmentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: `Create an environment`,
	Example: heredoc.Doc(`
	> cloud-platform environment create
	`),
	PreRun: upgradeIfNotLatest,
	RunE: func(cmd *cobra.Command, args []string) error {
		return environment.CreateTemplateNamespace(skipEnvCheck, answersFile)
	},
}

var environmentEcrCmd = &cobra.Command{
	Use:   "ecr",
	Short: `Add an ECR to a namespace`,
	Example: heredoc.Doc(`
	> cloud-platform environment ecr create
	`),
	PreRun: upgradeIfNotLatest,
}

var environmentPlanCmd = &cobra.Command{
	Use: "plan",
	Short: `Perform a terraform plan and kubectl apply -dry-run for a given namespace using either -namespace flag or the
	the namespace in the given PR Id/Number`,
	Long: `
	Perform a kubectl apply -dry-run and a terraform plan for a given namespace using either -namespace flag or the
	the namespace in the given PR Id/Number

	Along with the mandatory input flag, the below environments variables needs to be set
	TF_VAR_cluster_name - e.g. "cp-1902-02" to get the vpc details for some modules like rds, es
	TF_VAR_cluster_state_bucket - State where the cluster state is stored
	TF_VAR_cluster_state_key - folder name/state key inside the state bucket where cluster state is stored
	TF_VAR_github_owner - Github owner: ministryofjustice
	TF_VAR_github_token - Personal access token with repo scope to push github action secrets
	TF_VAR_kubernetes_cluster - Full name of the Cluster e.g. XXXXXX.gr7.eu-west2.eks.amazonaws.com
	PINGDOM_API_TOKEN - API Token to access pingdom
	PIPELINE_TERRAFORM_STATE_LOCK_TABLE - DynamoDB table where the state lock is stored
	PIPELINE_STATE_BUCKET - State bucket where the environments state is stored e.g cloud-platform-terraform-state
	PIPELINE_STATE_KEY_PREFIX - State key/ folder where the environments terraform state is stored e.g cloud-platform-environments
	PIPELINE_STATE_REGION - State region of the bucket e.g. eu-west-1
	PIPELINE_CLUSTER - Cluster name/folder inside namespaces/ in cloud-platform-environments
	PIPELINE_CLUSTER_STATE - Cluster name/folder inside the state bucket where the environments terraform state is stored. for "live" the state is stored under "live-1.cloud-platform.service..."
	`,
	Example: heredoc.Doc(`
	$ cloud-platform environment plan
	`),
	PreRun: upgradeIfNotLatest,
	Run: func(cmd *cobra.Command, args []string) {
		contextLogger := log.WithFields(log.Fields{"subcommand": "plan"})

		ghConfig := &github.GithubClientConfig{
			Repository: "cloud-platform-environments",
			Owner:      "ministryofjustice",
		}

		applier := &environment.Apply{
			Options:      &optFlags,
			GithubClient: github.NewGithubClient(ghConfig, optFlags.GithubToken),
		}

		err := applier.Plan()
		if err != nil {
			contextLogger.Fatal(err)
		}
	},
}

var environmentApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: `Perform a terraform apply and kubectl apply for a given namespace`,
	Long: `
	Perform a kubectl apply and a terraform apply for a given namespace using either -namespace flag or the
	the namespace in the given PR Id/Number

	Along with the mandatory input flag, the below environments variables needs to be set
	TF_VAR_cluster_name - e.g. "cp-1902-02" to get the vpc details for some modules like rds, es
	TF_VAR_cluster_state_bucket - State where the cluster state is stored
	TF_VAR_cluster_state_key - folder name/state key inside the state bucket where cluster state is stored
	TF_VAR_github_owner - Github owner: ministryofjustice
	TF_VAR_github_token - Personal access token with repo scope to push github action secrets
	TF_VAR_kubernetes_cluster - Full name of the Cluster e.g. XXXXXX.gr7.eu-west2.eks.amazonaws.com
	PINGDOM_API_TOKEN - API Token to access pingdom
	PIPELINE_TERRAFORM_STATE_LOCK_TABLE - DynamoDB table where the state lock is stored
	PIPELINE_STATE_BUCKET - State bucket where the environments state is stored e.g cloud-platform-terraform-state
	PIPELINE_STATE_KEY_PREFIX - State key/ folder where the environments terraform state is stored e.g cloud-platform-environments
	PIPELINE_STATE_REGION - State region of the bucket e.g. eu-west-1
	PIPELINE_CLUSTER - Cluster name/folder inside namespaces/ in cloud-platform-environments
	PIPELINE_CLUSTER_STATE - Cluster name/folder inside the state bucket where the environments terraform state is stored
	`,
	Example: heredoc.Doc(`
	$ cloud-platform environment apply -n <namespace>
	`),
	PreRun: upgradeIfNotLatest,
	Run: func(cmd *cobra.Command, args []string) {
		contextLogger := log.WithFields(log.Fields{"subcommand": "apply"})

		ghConfig := &github.GithubClientConfig{
			Repository: "cloud-platform-environments",
			Owner:      "ministryofjustice",
		}

		applier := &environment.Apply{
			Options:      &optFlags,
			GithubClient: github.NewGithubClient(ghConfig, optFlags.GithubToken),
		}

		if optFlags.AllNamespaces {
			err := applier.ApplyAll()
			if err != nil {
				contextLogger.Fatal(err)
			}
		} else {
			err := applier.Apply()
			if err != nil {
				contextLogger.Fatal(err)
			}
		}
	},
}

var environmentEcrCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  `Create "resources/ecr.tf" terraform file for an ECR`,
	PreRun: upgradeIfNotLatest,
	RunE:   environment.CreateTemplateEcr,
}

var environmentRdsCmd = &cobra.Command{
	Use:   "rds",
	Short: `Add an RDS instance to a namespace`,
	Example: heredoc.Doc(`
	> cloud-platform environment rds create
	`),
	PreRun: upgradeIfNotLatest,
}

var environmentRdsCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  `Create "resources/rds.tf" terraform file for an RDS instance`,
	PreRun: upgradeIfNotLatest,
	RunE:   environment.CreateTemplateRds,
}

var environmentS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: `Add a S3 bucket to a namespace`,
	Example: heredoc.Doc(`
	> cloud-platform environment s3 create
	`),
	PreRun: upgradeIfNotLatest,
}

var environmentS3CreateCmd = &cobra.Command{
	Use:    "create",
	Short:  `Create "resources/s3.tf" terraform file for a S3 bucket`,
	PreRun: upgradeIfNotLatest,
	RunE:   environment.CreateTemplateS3,
}

var environmentSvcCmd = &cobra.Command{
	Use:   "serviceaccount",
	Short: `Add a serviceaccount to a namespace`,
	Example: heredoc.Doc(`
	> cloud-platform environment serviceaccount
	`),
	PreRun: upgradeIfNotLatest,
}

var environmentSvcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: `Creates a serviceaccount in your chosen namespace`,
	Example: heredoc.Doc(`
	> cloud-platform environment serviceaccount create
	`),
	PreRun: upgradeIfNotLatest,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := environment.CreateTemplateServiceAccount(); err != nil {
			return err
		}

		return nil
	},
}

var environmentPrototypeCmd = &cobra.Command{
	Use:   "prototype",
	Short: `Create a gov.uk prototype kit site on the cloud platform`,
	Example: heredoc.Doc(`
	> cloud-platform environment prototype
	`),
	PreRun: upgradeIfNotLatest,
}

var environmentPrototypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: `Create an environment to host gov.uk prototype kit site on the cloud platform`,
	Long: `
Create a namespace folder and files in an existing prototype github repository to host a Gov.UK
Prototype Kit website on the Cloud Platform.

The namespace name should be your prototype github repository name:

  https://github.com/ministryofjustice/[repository name]
	`,
	Example: heredoc.Doc(`
	> cloud-platform environment prototype create
	`),
	PreRun: upgradeIfNotLatest,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := environment.CreateTemplatePrototype(); err != nil {
			return err
		}

		return nil
	},
}

var environmentBumpModuleCmd = &cobra.Command{
	Use:   "bump-module",
	Short: `Bump all specified module versions`,
	Example: heredoc.Doc(`
cloud-platform environments bump-module --module serviceaccount --module-version 1.1.1

Would bump all users serviceaccount modules in the environments repository to the specified version.
	`),
	PreRun: upgradeIfNotLatest,
	RunE: func(cmd *cobra.Command, args []string) error {
		if moduleVersion == "" || module == "" {
			return errors.New("--module and --module-version are required")
		}

		if err := environment.BumpModule(module, moduleVersion); err != nil {
			return err
		}
		return nil
	},
}
