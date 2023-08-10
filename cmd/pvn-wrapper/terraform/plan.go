package terraform

import (
	"context"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/prodvana/pvn-wrapper/result"
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
	Run: func(cmd *cobra.Command, args []string) {
		result.RunWrapper(func(ctx context.Context) (*result.ResultType, []result.OutputFileUpload, error) {
			const terraformOutFile = "plan.tfplan"
			planArgs := []string{"plan"}
			planArgs = append(planArgs, args...)
			planArgs = append(planArgs,
				"--detailed-exitcode",
				"--out",
				terraformOutFile,
			)
			execCmd := exec.CommandContext(ctx, terraformPath, planArgs...)
			res, err := result.RunCmd(execCmd)
			if err != nil {
				return res, nil, errors.Wrap(err, "plan command failed")
			}
			showCommand := exec.CommandContext(ctx, terraformPath, "show", terraformOutFile)
			output, err := showCommand.CombinedOutput()
			if err != nil {
				return res, nil, errors.Wrap(err, "show command failed")
			}
			return res, []result.OutputFileUpload{
				{
					Name: result.PlanOutput,
					Path: terraformOutFile,
				},
				{
					Name:    result.PlanExplanation,
					Content: output,
				},
			}, err
		})
	},
}

func init() {
	RootCmd.AddCommand(planCmd)
}
