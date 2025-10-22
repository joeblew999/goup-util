# Deprecated Packaging Files

These files are from the old packaging system and are **no longer used**.

## What's Here

### macOS Old Files

- **`macos-info.plist`** - Old hardcoded Info.plist
  - **Replaced by**: `pkg/packaging/templates/macos-info.plist.tmpl`
  - **Reason**: Now template-based with variables

- **`entitlements.plist`** - Old hardcoded entitlements
  - **Replaced by**: `pkg/packaging/templates/macos-entitlements.plist.tmpl`
  - **Reason**: Now template-based, embedded in Go

- **`build-macos.sh`** - Old bash packaging script
  - **Replaced by**: `pkg/packaging/macos.go` (pure Go)
  - **Reason**: Cross-platform, testable, better error handling

### Windows Old Commands

- **`msix_manifest.go`** - Old manifest generation command
  - **Replaced by**: `cmd/bundle.go` (unified command)
  - **Reason**: Integrated into unified packaging workflow

- **`msix_pack.go`** - Old MSIX packaging command
  - **Replaced by**: `cmd/bundle.go --create-msix` flag
  - **Reason**: Integrated into unified packaging workflow

## Migration

### Old Way (bash script)
```bash
bash pkg/packaging/build-macos.sh
```

### New Way (pure Go)
```bash
goup-util bundle macos examples/hybrid-dashboard
```

## Why These Files Are Kept

These files are preserved for **reference only** in case we need to compare implementations or restore old behavior. They are not used by the codebase.

## When to Delete

These files can be safely deleted once:
1. ✅ New packaging system is confirmed working
2. ✅ All tests pass on macOS and Windows
3. ✅ No users report issues with the new system

**Estimated deletion date**: After 1-2 release cycles (~Q1 2026)

## See Also

- `docs/PACKAGING.md` - New packaging system documentation
- `pkg/packaging/macos.go` - macOS bundle implementation
- `pkg/packaging/windows.go` - Windows bundle implementation
- `cmd/bundle.go` - Unified bundle command
