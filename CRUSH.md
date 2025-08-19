# CRUSH.md - goup-util Development Guide

## Build Commands
```bash
goup-util build [platform] [app-directory]    # Build examples using CLI's built-in build command
goup-util self build                          # Build goup-util itself
goup-util self upgrade                        # Upgrade goup-util to latest GitHub release
```

## Development Cycle
Complete 360° development: build → test → release → self-update. Use MCP and LSP tools for assistance throughout the entire cycle.

## Code Style
- **Imports**: stdlib → 3rd party → internal (separated by blank lines)
- **Naming**: PascalCase types/exports, camelCase vars/funcs, kebab-case files
- **Errors**: Wrap with `fmt.Errorf("...: %w", err)`, early returns preferred
- **Comments**: Package docs first, complete sentences for funcs, sparse inline
- **Formatting**: gofmt standard, blank lines between logical sections
- **CLI**: Use emojis for user feedback (✅📦❌)

## Project Structure
- `cmd/` - CLI commands (cobra)
- `pkg/` - Reusable packages
- `examples/` - Sample projects
- Go 1.25.0 required (see go.mod)