package output

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// SafeExecute wraps command execution with panic recovery
// Ensures JSON output even if command panics
func SafeExecute(command string, fn func() error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			base := &BaseResult{
				Command:   command,
				Version:   JSONSchemaVersion,
				Timestamp: time.Now().UTC(),
				Status:    StatusError,
				ExitCode:  ExitPanic,
				Error: &ErrorInfo{
					Message: fmt.Sprintf("panic: %v", r),
					Type:    ErrorTypePanic,
					Details: stack,
				},
			}
			printJSON(base)
			os.Exit(ExitPanic)
		}
	}()

	if err := fn(); err != nil {
		PrintError(command, err)
	}
}
