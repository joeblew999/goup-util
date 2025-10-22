package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Print outputs a Result as JSON to stdout
func Print(result Result, command string) {
	base := result.ToBaseResult(command)
	printJSON(base)
}

// PrintError outputs an error as JSON to stdout with non-zero exit code
func PrintError(command string, err error) {
	base := &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusError,
		ExitCode:  ExitError,
		Error: &ErrorInfo{
			Message: err.Error(),
			Type:    ErrorTypeExecution,
		},
	}
	printJSON(base)
	os.Exit(ExitError)
}

// PrintSuccess outputs a success result with optional data
func PrintSuccess(command string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	base := &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      jsonData,
	}
	printJSON(base)
}

// OK outputs a typed struct as JSON
// Usage: output.OK("self version", VersionResult{Version: "1.0.0", OS: "darwin"})
func OK[T any](command string, data T) {
	jsonData, _ := json.Marshal(data)
	base := &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      jsonData,
	}
	printJSON(base)
}

// Err outputs an error and exits
// Usage: output.Err("self version", err)
func Err(command string, err error) {
	PrintError(command, err)
}

// Run wraps any function to automatically output JSON on success or error
// Usage: output.Run("self build", func() (*BuildResult, error) { return &result, nil })
func Run[T any](command string, fn func() (*T, error)) {
	result, err := fn()
	if err != nil {
		PrintError(command, err)
		return
	}

	jsonData, _ := json.Marshal(result)
	base := &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      jsonData,
	}
	printJSON(base)
}

// printJSON safely encodes and outputs JSON
func printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		// Fallback: if encoding fails, output minimal error JSON
		fmt.Fprintf(os.Stdout, `{"command":"unknown","version":"%s","status":"%s","exit_code":%d,"error":{"message":"JSON encoding failed: %s"}}%s`, JSONSchemaVersion, StatusError, ExitPanic, err.Error(), "\n")
		os.Exit(ExitPanic)
	}
}
