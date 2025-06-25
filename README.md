# goup-util

Cross-platform SDK manager and build tool for Go applications.

**Developer-focused, idempotent, and DRY design.**

## Installation

```bash
go install github.com/joeblew999/goup-util@latest
```

## Build

Use [Task](https://taskfile.dev/) for development operations:

```bash
# Build the tool
task build

# List all available tasks
task --list
```

## Architecture

- **Idempotent operations**: Safe to run multiple times
- **DRY principles**: Centralized path management, shared utilities
- **Clean interfaces**: Service layer ready for future API use

## Documentation

The `docs/` folder contains end-user documentation and guides for using goup-util.
