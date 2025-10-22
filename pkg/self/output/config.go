package output

// JSON schema version
const JSONSchemaVersion = "1"

// Status values
const (
	StatusOK      = "ok"
	StatusWarning = "warning"
	StatusError   = "error"
)

// Error types
const (
	ErrorTypeExecution = "execution_error"
	ErrorTypePanic     = "panic"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
	ExitPanic   = 2
)

// Command names
const (
	CommandVersion   = "self version"
	CommandStatus    = "self status"
	CommandDoctor    = "self doctor"
	CommandBuild     = "self build"
	CommandSetup     = "self setup"
	CommandUninstall = "self uninstall"
	CommandTest      = "self test"
	CommandUpgrade   = "self upgrade"
	CommandRelease   = "self release"
)
