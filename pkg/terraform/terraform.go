// Package terraform implements methods and functions for running
// Terraform commands, such as terraform init/plan/apply.
package terraform

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	wsFailedToSelectRegexp = regexp.MustCompile(`Failed to select workspace`)
	wsDoesNotExistRegexp   = regexp.MustCompile(`workspace ".*" does not exist`)
)

// TerraformCLI is the client that wraps around terraform-exec
// to execute Terraform cli commands
type TerraformCLI struct {
	tf         terraformExec
	workingDir string
	workspace  string
	applyVars  []tfexec.ApplyOption
	planVars   []tfexec.PlanOption
	initVars   []tfexec.InitOption
	Redacted   bool
}

// TerraformCLIConfig configures the Terraform client
type TerraformCLIConfig struct {
	// ExecPath is the path to the Terraform executable.
	ExecPath string
	// WorkingDir is the path Terraform will execute in.
	WorkingDir string
	// Worspace is the Terraform workspace to use.
	Workspace string
	// ApplyVars allows you to group apply options passed to Terraform.
	ApplyVars []tfexec.ApplyOption
	// PlanVars allows you to group plan variables passed to Terraform.
	PlanVars []tfexec.PlanOption
	// InitVars allows you to group init variables passed to Terraform.
	InitVars []tfexec.InitOption
	// Version is the version of Terraform to use.
	Version string
	// Redacted is the flag to enable/disable redacting the terraform before printing output.
	Redacted bool
}

// NewTerraformCLI creates a terraform-exec client and configures and
// initializes a new Terraform client
func NewTerraformCLI(config *TerraformCLIConfig) (*TerraformCLI, error) {
	if config == nil {
		return nil, errors.New("TerraformCLIConfig cannot be nil - no meaningful default values")
	}

	if config.Version == "" {
		return nil, errors.New("version cannot be empty")
	}

	if config.ExecPath == "" {
		i := install.NewInstaller()
		v := version.Must(version.NewVersion(config.Version))

		execPath, err := i.Ensure(context.TODO(), []src.Source{
			&fs.ExactVersion{
				Product: product.Terraform,
				Version: v,
			},
			&releases.ExactVersion{
				Product: product.Terraform,
				Version: v,
			},
		})
		if err != nil {
			return nil, err
		}

		config.ExecPath = execPath

		defer func() {
			if err := i.Remove(context.TODO()); err != nil {
				return
			}
		}()
	}

	tf, err := tfexec.NewTerraform(config.WorkingDir, config.ExecPath)
	if err != nil {
		return nil, err
	}

	client := &TerraformCLI{
		tf:         tf,
		workingDir: config.WorkingDir,
		workspace:  config.Workspace,
		applyVars:  config.ApplyVars,
		planVars:   config.PlanVars,
		initVars:   config.InitVars,
		Redacted:   config.Redacted,
	}

	return client, nil
}

// Init initializes by executing the cli command `terraform init` and
// `terraform workspace new <name>`
func (t *TerraformCLI) Init(ctx context.Context, w io.Writer) error {
	var wsCreated bool

	t.tf.SetStdout(w)
	t.tf.SetStderr(w)
	// This is special handling for when the workspace has been detected in
	// .terraform/environment with a non-existing state. This case is common
	// when the state for the workspace has been deleted.
	// https://github.com/hashicorp/terraform/issues/21393
TF_INIT_AGAIN:
	if err := t.tf.Init(ctx); err != nil {
		var wsErr *tfexec.ErrNoWorkspace
		matchedFailedToSelect := wsFailedToSelectRegexp.MatchString(err.Error())
		matchedDoesNotExist := wsDoesNotExistRegexp.MatchString(err.Error())
		if matchedFailedToSelect || matchedDoesNotExist || errors.As(err, &wsErr) {
			fmt.Println("workspace was detected without state, " +
				"creating new workspace and attempting Terraform init again")
			if err := t.tf.WorkspaceNew(ctx, t.workspace); err != nil {
				return err
			}

			if !wsCreated {
				wsCreated = true
				goto TF_INIT_AGAIN
			}
		}
		return err
	}

	if !wsCreated {
		err := t.tf.WorkspaceNew(ctx, t.workspace)
		if err != nil {
			var wsErr *tfexec.ErrWorkspaceExists
			if !errors.As(err, &wsErr) {
				return err
			}
		}
	}

	if err := t.tf.WorkspaceSelect(ctx, t.workspace); err != nil {
		return err
	}

	return nil
}

// Apply executes the cli command `terraform apply` for a given workspace
func (t *TerraformCLI) Apply(ctx context.Context, w io.Writer) error {

	t.tf.SetStdout(w)
	t.tf.SetStderr(w)

	if err := t.tf.Apply(ctx, t.applyVars...); err != nil {
		return err
	}

	return nil
}

// Plan executes the cli command `terraform plan` for a given workspace
func (t *TerraformCLI) Plan(ctx context.Context, w io.Writer) (bool, error) {

	t.tf.SetStdout(w)
	t.tf.SetStderr(w)

	diff, err := t.tf.Plan(ctx, t.planVars...)

	if err != nil {
		return false, err
	}

	return diff, nil
}

// Plan executes the cli command `terraform plan` for a given workspace
func (t *TerraformCLI) Output(ctx context.Context, w io.Writer) (map[string]tfexec.OutputMeta, error) {
	// Often times, the output is not needed, so the caller can specify a null writer to ignore.
	t.tf.SetStdout(w)
	return t.tf.Output(ctx)
}
