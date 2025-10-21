# Claude Instructions for goup-util

## Project Overview

**goup-util** is a cross-platform SDK manager and build tool for Go applications, specifically designed for building Gio UI applications for Android, iOS, macOS, and Windows platforms.

### Key Principles
- **Idempotent operations**: All operations are safe to run multiple times
- **DRY (Don't Repeat Yourself)**: Centralized path management and shared utilities
- **Developer-focused**: Clean CLI interface with clear commands
- **Cross-platform**: Supports macOS, Linux, Windows, Android, and iOS

## Project Structure

```
goup-util/
├── cmd/                    # CLI commands (Cobra-based)
│   ├── root.go            # Root command
│   ├── build.go           # Build Gio apps for platforms
│   ├── install.go         # Install SDKs
│   ├── self.go            # Build/update goup-util itself
│   ├── icons.go           # Generate platform icons
│   ├── package.go         # Package apps for distribution
│   ├── workspace.go       # Manage Go workspaces
│   └── ...
├── pkg/                   # Shared packages
│   ├── config/           # Configuration management
│   ├── installer/        # SDK installation logic
│   ├── icons/            # Icon generation
│   ├── workspace/        # Go workspace utilities
│   └── ...
├── examples/             # Example Gio applications
│   ├── gio-basic/
│   ├── gio-plugin-hyperlink/
│   └── gio-plugin-webviewer/
├── docs/                 # End-user documentation
└── main.go              # Entry point

```

## Development Workflow

### Building the Tool

```bash
# Build goup-util itself
go run . self build

# Build for all platforms
go run . self build --all

# Run tests
go test ./...

# Run integration tests
go test -v integration_test.go
```

### Testing Changes

```bash
# Use 'go run .' to test changes without building
go run . build macos examples/gio-basic
go run . install android-sdk
go run . icons examples/gio-basic
```

## Key Commands to Understand

- `build` - Build Gio applications for different platforms (macos, windows, android, ios)
- `install` - Install SDKs (Android SDK, NDK, etc.)
- `self build` - Build goup-util binaries for distribution
- `icons` - Generate platform-specific icons from source images
- `package` - Package built apps for distribution
- `workspace` - Manage Go workspace files
- `gitignore` - Manage .gitignore templates

## Important Files

- `cmd/*.go` - All CLI command implementations
- `pkg/config/` - Config file handling and directory management
- `pkg/installer/` - SDK installation logic
- `go.mod` - Dependencies (cobra, progressbar, icns, etc.)
- `.gitignore` - Build binaries are excluded (goup-util*)

## Common Tasks

### Adding a New Command

1. Create `cmd/newcommand.go`
2. Use Cobra patterns from existing commands
3. Add to root command in `cmd/root.go`
4. Test with `go run . newcommand`

### Modifying Build Logic

- See `cmd/build.go` for platform-specific build commands
- Build logic uses idempotent patterns
- Platform detection and SDK path resolution in `pkg/`

### Working with Icons

- Icon generation in `pkg/icons/`
- Supports PNG source → platform formats (icns, ico, Android drawables)
- Test with example projects

## Dependencies

Key external packages:
- `github.com/spf13/cobra` - CLI framework
- `github.com/schollz/progressbar/v3` - Progress display
- `github.com/JackMordaunt/icns` - macOS icon generation
- Platform-specific SDK tools (Android SDK, Xcode, etc.)

## Testing Guidelines

- Test commands using `go run .` before building
- Use example projects in `examples/` for integration testing
- Verify idempotency (running twice should produce same result)
- Test on target platforms when modifying build logic

## CI/CD

- GitHub Actions in `.github/workflows/`
- `build.yml` - Main CI pipeline using `go run . self build`
- Builds binaries for all platforms
- Artifacts uploaded as releases

## Future Plans (See TODO.md)

- UTM integration for Windows VM testing
- Winget package management for Windows dependencies
- Automated testing infrastructure

## Tips for Claude

1. **Always test with `go run .`** - Don't build binaries during development
2. **Maintain idempotency** - Operations should be safe to run multiple times
3. **Follow existing patterns** - Look at similar commands for consistency
4. **Update docs/** - Keep end-user docs in sync with code changes
5. **Check .gitignore** - Don't commit build binaries (goup-util*)
6. **Use examples/** - Test changes with the example Gio projects
7. **Cross-platform awareness** - Code runs on macOS, Linux, and Windows

## Common Debugging

```bash
# Check configuration
go run . config

# List available SDKs
go run . list

# Verbose output (add -v flag if available)
go run . build macos examples/gio-basic -v

# Check Go workspace
go run . workspace list
```

## Code Style

- Follow standard Go conventions
- Use `cobra` for CLI structure
- Error handling with clear messages
- Progress bars for long operations
- Idempotent file operations
