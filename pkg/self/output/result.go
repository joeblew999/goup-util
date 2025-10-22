package output

import (
	"encoding/json"
	"fmt"
	"time"
)

// BaseResult is the universal structure for all command outputs
type BaseResult struct {
	Command   string          `json:"command"`
	Version   string          `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Status    string          `json:"status"` // "ok", "warning", "error"
	ExitCode  int             `json:"exit_code"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     *ErrorInfo      `json:"error,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}

// Result is the interface all command results must implement
type Result interface {
	ToBaseResult(command string) *BaseResult
}

// VersionResult represents version command output
type VersionResult struct {
	Version  string `json:"version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Location string `json:"location,omitempty"`
}

func (v VersionResult) ToBaseResult(command string) *BaseResult {
	data, _ := json.Marshal(v)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// StatusResult represents status command output
type StatusResult struct {
	Installed       bool   `json:"installed"`
	CurrentVersion  string `json:"current_version,omitempty"`
	LatestVersion   string `json:"latest_version,omitempty"`
	UpdateAvailable bool   `json:"update_available"`
	Location        string `json:"location,omitempty"`
}

func (s StatusResult) ToBaseResult(command string) *BaseResult {
	status := StatusOK
	if !s.Installed {
		status = StatusWarning
	}
	data, _ := json.Marshal(s)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    status,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// DoctorResult represents doctor command output
type DoctorResult struct {
	Installations []InstallationInfo `json:"installations"`
	Issues        []string           `json:"issues,omitempty"`
	Suggestions   []string           `json:"suggestions,omitempty"`
}

type InstallationInfo struct {
	Path     string `json:"path"`
	Active   bool   `json:"active"`
	Shadowed bool   `json:"shadowed"`
}

func (d DoctorResult) ToBaseResult(command string) *BaseResult {
	status := StatusOK
	exitCode := ExitSuccess

	if len(d.Installations) == 0 {
		status = StatusError
		exitCode = ExitError
	} else if len(d.Installations) > 1 || len(d.Issues) > 0 {
		status = StatusWarning
	}

	data, _ := json.Marshal(d)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    status,
		ExitCode:  exitCode,
		Data:      data,
	}
}

// BuildResult represents build command output
type BuildResult struct {
	Binaries         []string `json:"binaries"`
	ScriptsGenerated bool     `json:"scripts_generated"`
	OutputDir        string   `json:"output_dir"`
	LocalMode        bool     `json:"local_mode"`
}

func (b BuildResult) ToBaseResult(command string) *BaseResult {
	data, _ := json.Marshal(b)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// SetupResult represents setup command output
type SetupResult struct {
	Installed      bool   `json:"installed"`
	Location       string `json:"location"`
	InPath         bool   `json:"in_path"`
	DependenciesOK bool   `json:"dependencies_ok"`
}

func (s SetupResult) ToBaseResult(command string) *BaseResult {
	status := StatusOK
	if !s.DependenciesOK {
		status = StatusWarning
	}
	data, _ := json.Marshal(s)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    status,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// UninstallResult represents uninstall command output
type UninstallResult struct {
	Removed []string `json:"removed"`
	Failed  []string `json:"failed,omitempty"`
}

func (u UninstallResult) ToBaseResult(command string) *BaseResult {
	status := StatusOK
	exitCode := ExitSuccess
	if len(u.Failed) > 0 {
		status = StatusWarning
		exitCode = ExitError
	}
	if len(u.Removed) == 0 && len(u.Failed) == 0 {
		status = StatusWarning
	}
	data, _ := json.Marshal(u)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    status,
		ExitCode:  exitCode,
		Data:      data,
	}
}

// TestResult represents test command output
type TestResult struct {
	Phase  string   `json:"phase"`
	Passed bool     `json:"passed"`
	Steps  []string `json:"steps"`
	Errors []string `json:"errors,omitempty"`
}

func (t TestResult) ToBaseResult(command string) *BaseResult {
	status := StatusOK
	exitCode := ExitSuccess
	if !t.Passed {
		status = StatusError
		exitCode = ExitError
	}
	data, _ := json.Marshal(t)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    status,
		ExitCode:  exitCode,
		Data:      data,
	}
}

// Parse*Data helper methods for bidirectional JSON parsing

// ParseVersionData parses Data field as VersionResult
func (b *BaseResult) ParseVersionData() (*VersionResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result VersionResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseStatusData parses Data field as StatusResult
func (b *BaseResult) ParseStatusData() (*StatusResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result StatusResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseDoctorData parses Data field as DoctorResult
func (b *BaseResult) ParseDoctorData() (*DoctorResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result DoctorResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseBuildData parses Data field as BuildResult
func (b *BaseResult) ParseBuildData() (*BuildResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result BuildResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseSetupData parses Data field as SetupResult
func (b *BaseResult) ParseSetupData() (*SetupResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result SetupResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseUninstallData parses Data field as UninstallResult
func (b *BaseResult) ParseUninstallData() (*UninstallResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result UninstallResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseTestData parses Data field as TestResult
func (b *BaseResult) ParseTestData() (*TestResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result TestResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// UpgradeResult represents upgrade command output
type UpgradeResult struct {
	PreviousVersion string `json:"previous_version,omitempty"`
	NewVersion      string `json:"new_version"`
	Downloaded      bool   `json:"downloaded"`
	Installed       bool   `json:"installed"`
	Location        string `json:"location"`
}

func (u UpgradeResult) ToBaseResult(command string) *BaseResult {
	data, _ := json.Marshal(u)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// ReleaseResult represents release command output
type ReleaseResult struct {
	Version      string   `json:"version"`
	TestsPassed  bool     `json:"tests_passed"`
	Built        bool     `json:"built"`
	Tagged       bool     `json:"tagged"`
	Pushed       bool     `json:"pushed"`
	Binaries     []string `json:"binaries,omitempty"`
}

func (r ReleaseResult) ToBaseResult(command string) *BaseResult {
	data, _ := json.Marshal(r)
	return &BaseResult{
		Command:   command,
		Version:   JSONSchemaVersion,
		Timestamp: time.Now().UTC(),
		Status:    StatusOK,
		ExitCode:  ExitSuccess,
		Data:      data,
	}
}

// ParseUpgradeData parses Data field as UpgradeResult
func (b *BaseResult) ParseUpgradeData() (*UpgradeResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result UpgradeResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}

// ParseReleaseData parses Data field as ReleaseResult
func (b *BaseResult) ParseReleaseData() (*ReleaseResult, error) {
	if b.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}
	var result ReleaseResult
	err := json.Unmarshal(b.Data, &result)
	return &result, err
}
