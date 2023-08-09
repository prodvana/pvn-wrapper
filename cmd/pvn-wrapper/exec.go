package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

const (
	PvnWrapperVersion = "0.0.2"
)

type ResultType struct {
	ExitCode         int    `json:"exit_code"`  // Exit code of wrapped process. -1 if process failed to execute.
	ExecError        string `json:"exec_error"` // Internal error when trying to execute wrapped process.
	Stdout           []byte `json:"stdout"`
	Stderr           []byte `json:"stderr"`
	Version          string `json:"version"`     // Wrapper version.
	StartTimestampNs int64  `json:"start_ts_ns"` // Timestamp when the process began executing, in ns.
	DurationNs       int64  `json:"duration_ns"` // Total execution duration of the process, in ns.
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command then wrap its output in a format that Prodvana understands.",
	Long: `Execute a command then wrap its output in a format that Prodvana understands.
The exit code matches the exit code of the underlying binary being executed.

pvn-wrapper exec my-binary --my-flag=value my-args ...
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startTs := time.Now()

		execCmd := exec.Command(args[0], args[1:]...)
		execCmd.Env = os.Environ()

		// TODO: Limit stdout/stderr to a reasonable size while preserving useful error context.
		// Kubernetes output is usually limited to 10MB.
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		execCmd.Stdout = stdout
		execCmd.Stderr = stderr

		var result ResultType

		err := execCmd.Run()
		duration := time.Since(startTs)

		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				result.ExitCode = exitErr.ExitCode()
			} else {
				result.ExecError = err.Error()
				result.ExitCode = -1
			}
		}

		result.Stdout = stdout.Bytes()
		result.Stderr = stderr.Bytes()
		result.Version = PvnWrapperVersion
		result.StartTimestampNs = startTs.UnixNano()
		result.DurationNs = duration.Nanoseconds()

		err = json.NewEncoder(os.Stdout).Encode(&result)
		if err != nil {
			// If something went wrong during encode/write to stdout, indicate that in stderr and exit non-zero.
			log.Fatal(err)
		}

		// If the wrapped process fails, make sure this process has a non-zero exit code.
		// This is to maintain compatibility with existing task execution infrastructure.
		// Once we enforce the use of this wrapper, we can safely exit 0 here.
		os.Exit(result.ExitCode)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
