package main

import (
	"bytes"
	"errors"
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
		result.RunWrapper(func() (*result.ResultType, error) {
			execCmd := exec.Command(args[0], args[1:]...)
			execCmd.Env = os.Environ()

			// TODO: Limit stdout/stderr to a reasonable size while preserving useful error context.
			// Kubernetes output is usually limited to 10MB.
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			execCmd.Stdout = stdout
			execCmd.Stderr = stderr

			var result result.ResultType

			err := execCmd.Run()

			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					result.ExitCode = exitErr.ExitCode()
				} else {
					return nil, err
				}
			}

			result.Stdout = stdout.Bytes()
			result.Stderr = stderr.Bytes()
			return &result, nil
		})
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
