package terraform

import (
	"bytes"
	"os/exec"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:     "plan",
	Short:   "terraform plan wrapper",
	Aliases: []string{"tf"},
	Long: `terraform plan wrapper.

Takes all the same input that terraform plan would, but handles uploading plan file to Prodvana and exiting with 0, 1, or 2.

0 - No changes detected
1 - Unknown error
2 - Changes detected

pvn-wrapper terraform plan ...

To pass flags to terraform plan, use --

pvn-wrapper terraform plan -- --refresh=false

pvn-wrapper will always pass --detailed-exitcode and --out.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var planArgs []string
		planArgs = append(planArgs, args...)
		planArgs = append(planArgs,
			"--detailed-exitcode",
			"--out",
		)
		execCmd := exec.Command(terraformPath, planArgs...)

		// TODO: Limit stdout/stderr to a reasonable size while preserving useful error context.
		// Kubernetes output is usually limited to 10MB.
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		execCmd.Stdout = stdout
		execCmd.Stderr = stderr

	},
}

func init() {
	RootCmd.AddCommand(planCmd)
}
