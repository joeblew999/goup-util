// Package schema provides the SINGLE SOURCE OF TRUTH for all type definitions.
// Used by: CLI (Cobra), MCP tools, OpenAPI specs, LLM function schemas.
//
// Define types here with jsonschema tags → use everywhere.
package schema

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// =============================================================================
// PLATFORM - Build target platforms
// =============================================================================

// Platform represents a build target platform.
// Used by: build, run, bundle, package commands and MCP tools.
type Platform string

const (
	PlatformMacOS   Platform = "macos"
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
	PlatformWeb     Platform = "web"
)

// Platforms is the list of all valid platforms (for CLI completion)
var Platforms = []string{
	string(PlatformMacOS),
	string(PlatformIOS),
	string(PlatformAndroid),
	string(PlatformWindows),
	string(PlatformLinux),
	string(PlatformWeb),
}

// PlatformDescriptions for CLI tab completion
var PlatformDescriptions = map[string]string{
	string(PlatformMacOS):   "macOS desktop application",
	string(PlatformIOS):     "iOS mobile application",
	string(PlatformAndroid): "Android mobile application",
	string(PlatformWindows): "Windows desktop application",
	string(PlatformLinux):   "Linux desktop application",
	string(PlatformWeb):     "Web application (WASM)",
}

// ValidPlatform checks if a string is a valid platform
func ValidPlatform(s string) bool {
	for _, p := range Platforms {
		if p == s {
			return true
		}
	}
	return false
}

// =============================================================================
// VM STATUS - UTM virtual machine states
// =============================================================================

// VMStatus represents the state of a virtual machine.
type VMStatus string

const (
	VMStatusRunning   VMStatus = "running"
	VMStatusStopped   VMStatus = "stopped"
	VMStatusSuspended VMStatus = "suspended"
	VMStatusStarting  VMStatus = "starting"
)

// VMStatuses is the list of all valid VM statuses (for schema generation)
var VMStatuses = []VMStatus{
	VMStatusRunning,
	VMStatusStopped,
	VMStatusSuspended,
	VMStatusStarting,
}

// =============================================================================
// BUILD - Input/Output types for build operations
// =============================================================================

// BuildInput defines parameters for building an application.
// Tags: json (field name), jsonschema (description for MCP/API)
// Required: fields WITHOUT omitempty are required
type BuildInput struct {
	Platform Platform `json:"platform" jsonschema:"Target platform to build for"`
	AppDir   string   `json:"app_dir" jsonschema:"Path to application directory"`
	Force    bool     `json:"force,omitempty" jsonschema:"Force rebuild even if up-to-date"`
	Check    bool     `json:"check,omitempty" jsonschema:"Check if rebuild needed without building"`
	Schemes  string   `json:"schemes,omitempty" jsonschema:"Deep linking URI schemes (comma-separated)"`
	Queries  string   `json:"queries,omitempty" jsonschema:"Android app package queries (comma-separated)"`
	SignKey  string   `json:"signkey,omitempty" jsonschema:"Signing key path or name"`
}

// BuildOutput is the result of a build operation.
type BuildOutput struct {
	Success    bool   `json:"success"`
	OutputPath string `json:"output_path" jsonschema:"Path to built artifact"`
	Cached     bool   `json:"cached" jsonschema:"Whether result was from cache"`
	Duration   string `json:"duration,omitempty" jsonschema:"Build duration"`
}

// =============================================================================
// UTM - Input/Output types for VM operations
// =============================================================================

// UTMStartInput defines parameters for starting a VM.
type UTMStartInput struct {
	VMName string `json:"vm_name" jsonschema:"Name of the VM to start"`
}

// UTMStopInput defines parameters for stopping a VM.
type UTMStopInput struct {
	VMName string `json:"vm_name" jsonschema:"Name of the VM to stop"`
	Force  bool   `json:"force,omitempty" jsonschema:"Force stop without graceful shutdown"`
}

// UTMStatusInput defines parameters for checking VM status.
type UTMStatusInput struct {
	VMName string `json:"vm_name" jsonschema:"Name of the VM to check"`
}

// UTMOutput is the result of a VM operation.
type UTMOutput struct {
	Success bool     `json:"success"`
	VMName  string   `json:"vm_name"`
	Status  VMStatus `json:"status" jsonschema:"Current VM status"`
	Message string   `json:"message,omitempty"`
}

// VMInfo describes a virtual machine.
type VMInfo struct {
	Name   string   `json:"name"`
	Status VMStatus `json:"status" jsonschema:"Current VM status"`
	OS     string   `json:"os,omitempty" jsonschema:"Operating system"`
	Memory string   `json:"memory,omitempty" jsonschema:"Allocated memory"`
}

// UTMListOutput is the result of listing VMs.
type UTMListOutput struct {
	VMs []VMInfo `json:"vms"`
}

// =============================================================================
// ICONS - Input/Output types for icon generation
// =============================================================================

// IconsInput defines parameters for generating icons.
type IconsInput struct {
	AppDir string `json:"app_dir" jsonschema:"Path to application directory"`
}

// IconsOutput is the result of icon generation.
type IconsOutput struct {
	Success   bool     `json:"success"`
	Generated []string `json:"generated" jsonschema:"List of generated icon files"`
}

// =============================================================================
// INSTALL - Input/Output types for SDK installation
// =============================================================================

// InstallInput defines parameters for installing an SDK.
type InstallInput struct {
	SDKName string `json:"sdk_name" jsonschema:"Name of the SDK to install"`
	Force   bool   `json:"force,omitempty" jsonschema:"Force reinstall even if present"`
}

// InstallOutput is the result of SDK installation.
type InstallOutput struct {
	Success     bool   `json:"success"`
	SDKName     string `json:"sdk_name"`
	Version     string `json:"version" jsonschema:"Installed version"`
	InstallPath string `json:"install_path" jsonschema:"Installation path"`
}

// =============================================================================
// PACKAGE - Input/Output types for packaging
// =============================================================================

// PackageInput defines parameters for packaging an application.
type PackageInput struct {
	Platform Platform `json:"platform" jsonschema:"Target platform for package"`
	AppDir   string   `json:"app_dir" jsonschema:"Path to application directory"`
	BundleID string   `json:"bundle_id,omitempty" jsonschema:"Bundle identifier (e.g., com.company.app)"`
	Version  string   `json:"version,omitempty" jsonschema:"Application version"`
}

// PackageOutput is the result of packaging.
type PackageOutput struct {
	Success     bool   `json:"success"`
	PackagePath string `json:"package_path" jsonschema:"Path to created package"`
	Format      string `json:"format" jsonschema:"Package format (tar.gz, zip, apk, etc.)"`
}

// =============================================================================
// SCHEMA OPTIONS - For JSON Schema generation with enum support
// =============================================================================

// SchemaOptions returns jsonschema.ForOptions with TypeSchemas for all enums.
// Use this when generating JSON schemas to get proper enum values:
//
//	schema, err := jsonschema.For[BuildInput](schema.SchemaOptions())
//
// This enables the SSOT pattern: define types once, generate schemas with
// proper enum values for MCP tools, OpenAPI specs, and LLM function schemas.
func SchemaOptions() *jsonschema.ForOptions {
	// Convert Platform constants to []any for jsonschema
	platformEnums := make([]any, len(Platforms))
	for i, p := range Platforms {
		platformEnums[i] = p
	}

	// Convert VMStatus constants to []any for jsonschema
	vmStatusEnums := make([]any, len(VMStatuses))
	for i, s := range VMStatuses {
		vmStatusEnums[i] = s
	}

	return &jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeFor[Platform](): {
				Type: "string",
				Enum: platformEnums,
			},
			reflect.TypeFor[VMStatus](): {
				Type: "string",
				Enum: vmStatusEnums,
			},
		},
	}
}
