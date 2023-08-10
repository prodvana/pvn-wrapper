package main

import (
	"context"
	"os"
	"os/exec"

	"github.com/prodvana/pvn-wrapper/result"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command then wrap its output in a format that Prodvana understands.",
	Long: `Execute a command then wrap its output in a format that Prodvana understands.
The exit code matches the exit code of the underlying binary being executed.

pvn-wrapper exec my-binary --my-flag=value my-args ...
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result.RunWrapper(func(ctx context.Context) (*result.ResultType, []result.OutputFileUpload, error) {
			execCmd := exec.CommandContext(ctx, args[0], args[1:]...)
			execCmd.Env = os.Environ()

			res, err := result.RunCmd(execCmd)
			return res, nil, err
		})
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
