package result

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type ResultType struct {
	ExitCode         int    `json:"exit_code"`  // Exit code of wrapped process. -1 if process failed to execute.
	ExecError        string `json:"exec_error"` // Internal error when trying to execute wrapped process.
	Stdout           []byte `json:"stdout"`
	Stderr           []byte `json:"stderr"`
	Version          string `json:"version"`     // Wrapper version.
	StartTimestampNs int64  `json:"start_ts_ns"` // Timestamp when the process began executing, in ns.
	DurationNs       int64  `json:"duration_ns"` // Total execution duration of the process, in ns.
	Files            []File
}

type File struct {
	AbsPath       string `json:"abs_path"`
	ContentBlobId string `json:"content_blob_id"`
}

const (
	PvnWrapperVersion = "0.0.2"
)

func RunWrapper(run func() (*ResultType, error)) {
	startTs := time.Now()
	result, err := run()
	duration := time.Since(startTs)
	if err != nil {
		result := &ResultType{}
		result.ExecError = err.Error()
		result.ExitCode = -1
	}
	result.StartTimestampNs = result.StartTimestampNs
	result.DurationNs = duration.Nanoseconds()
	result.Version = PvnWrapperVersion

	err = json.NewEncoder(os.Stdout).Encode(&result)
	if err != nil {
		// If something went wrong during encode/write to stdout, indicate that in stderr and exit non-zero.
		log.Fatal(err)
	}

	// If the wrapped process fails, make sure this process has a non-zero exit code.
	// This is to maintain compatibility with existing task execution infrastructure.
	// Once we enforce the use of this wrapper, we can safely exit 0 here.
	os.Exit(result.ExitCode)
}
