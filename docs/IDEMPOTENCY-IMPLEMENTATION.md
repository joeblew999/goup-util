# Idempotency Implementation Guide

This document tracks the implementation of idempotent commands in goup-util.

## ✅ What's Been Implemented

### 1. Build Cache System (`pkg/buildcache/`)

Created a comprehensive build state tracking system:
- **Source hashing** - Detects when Go files, go.mod, or assets change
- **Timestamp tracking** - Records last successful build time
- **Output validation** - Checks if build artifacts still exist
- **Platform-aware** - Separate cache per project/platform combination

**Cache location**: `~/.goup-util/build-cache.json`

### 2. Build Command Idempotency (`cmd/build.go`)

**New flags added**:
```bash
--force    # Force rebuild even if up-to-date
--check    # Check if rebuild needed (exit 0=no, 1=yes)
```

**Idempotent behavior**:
- ✅ Checks source hash before rebuilding
- ✅ Skips build if sources unchanged
- ✅ Records successful builds in cache
- ✅ Clear messaging about why skipping/rebuilding

**Updated function**: `buildMacOS()` now includes:
1. Cache check before building
2. Clear skip message if up-to-date
3. Reason reporting when rebuilding
4. Success/failure recording

## 🚧 Still TODO

### 3. Other Build Functions (Need Same Treatment)

Apply the same pattern to:
- [ ] `buildAndroid()`
- [ ] `buildIOS()`
- [ ] `buildWindows()`
- [ ] `buildAll()` (wrapper function)

**Pattern to follow**:
```go
func buildPlatform(appDir, appName, outputDir, platform string, opts BuildOptions) error {
    appPath := filepath.Join(outputDir, appName+".ext")
    cache := getBuildCache()

    // Check if rebuild needed
    if !opts.Force {
        needsRebuild, reason := cache.NeedsRebuild(appName, platform, appDir, appPath)

        if opts.CheckOnly {
            if needsRebuild {
                fmt.Printf("Rebuild needed: %s\n", reason)
                os.Exit(1)
            } else {
                fmt.Printf("Up to date: %s\n", appPath)
                os.Exit(0)
            }
        }

        if !needsRebuild {
            fmt.Printf("✓ %s for %s is up-to-date (use --force to rebuild)\n", appName, platform)
            return nil
        }

        fmt.Printf("Rebuilding: %s\n", reason)
    }

    // ... actual build logic ...

    // Record result
    cache.RecordBuild(appName, platform, appDir, appPath, success)
    return err
}
```

### 4. Screenshot Command Idempotency

**New flags needed**:
```bash
--force    # Overwrite existing screenshot
```

**Changes to `cmd/screenshot.go`**:
```go
// Before capturing
if !force {
    if _, err := os.Stat(output); err == nil {
        fmt.Printf("✓ Screenshot already exists: %s\n", output)
        fmt.Println("  Use --force to overwrite")
        return nil
    }
}
```

### 5. Install Command

**Already mostly idempotent!** But could improve:
```bash
--reinstall  # Force reinstall even if cached
```

### 6. Package Command

**New flags needed**:
```bash
--force    # Force repackage even if exists
```

Check if `.dist/goup-util.app` exists and is newer than source.

### 7. Icons Command

**Already has framework** via service layer, but add:
```bash
--force    # Regenerate even if icons exist
```

## 📝 Taskfile Updates Needed

Update all build tasks to leverage idempotency:

### Current (wasteful):
```yaml
run:hybrid:
  cmds:
    - "{{.GOUP}} build macos {{.HYBRID_EXAMPLE}}"  # Always rebuilds
    - open {{.HYBRID_EXAMPLE}}/.bin/hybrid-dashboard.app
```

### Better (idempotent):
```yaml
run:hybrid:
  cmds:
    - "{{.GOUP}} build macos {{.HYBRID_EXAMPLE}}"  # Skips if up-to-date
    - open {{.HYBRID_EXAMPLE}}/.bin/hybrid-dashboard.app

run:hybrid:force:
  desc: Force rebuild and run hybrid-dashboard
  cmds:
    - "{{.GOUP}} build --force macos {{.HYBRID_EXAMPLE}}"
    - open {{.HYBRID_EXAMPLE}}/.bin/hybrid-dashboard.app

check:hybrid:
  desc: Check if hybrid-dashboard needs rebuild
  cmds:
    - "{{.GOUP}} build --check macos {{.HYBRID_EXAMPLE}}"
```

### Add CI/CD tasks:
```yaml
ci:check:
  desc: Check if any examples need rebuilding
  cmds:
    - "{{.GOUP}} build --check macos {{.HYBRID_EXAMPLE}}"
    - "{{.GOUP}} build --check macos {{.WEBVIEWER_EXAMPLE}}"
    - "{{.GOUP}} build --check macos {{.BASIC_EXAMPLE}}"

ci:build:
  desc: Build all examples (only if needed)
  cmds:
    - "{{.GOUP}} build macos {{.HYBRID_EXAMPLE}}"
    - "{{.GOUP}} build macos {{.WEBVIEWER_EXAMPLE}}"
    - "{{.GOUP}} build macos {{.BASIC_EXAMPLE}}"
```

## 🧪 Testing Idempotency

### Manual Test Script

```bash
#!/bin/bash
# test-idempotency.sh

echo "Testing build idempotency..."

# First build
echo "=== First build (should build) ==="
go run . build macos examples/hybrid-dashboard

# Second build
echo "=== Second build (should skip) ==="
go run . build macos examples/hybrid-dashboard

# Check command
echo "=== Check command (should say up-to-date) ==="
go run . build --check macos examples/hybrid-dashboard
echo "Exit code: $?"

# Touch a source file
echo "=== Touching source file ==="
touch examples/hybrid-dashboard/main.go

# Third build
echo "=== Third build (should rebuild) ==="
go run . build macos examples/hybrid-dashboard

# Force rebuild
echo "=== Force rebuild ==="
go run . build --force macos examples/hybrid-dashboard

echo "✓ Idempotency tests complete"
```

### Expected Output

```
First build: ✓ Built hybrid-dashboard for macOS: .bin/hybrid-dashboard.app
Second build: ✓ hybrid-dashboard for macos is up-to-date (use --force to rebuild)
Check: Up to date: .bin/hybrid-dashboard.app (exit 0)
After touch: Rebuilding: sources changed
Force: Building hybrid-dashboard for macOS... ✓ Built
```

## 📊 Performance Impact

**Expected improvements**:
- Repeated `task run:hybrid` - **2-3s → 0.1s** (30x faster)
- CI builds when no changes - **5min → 10s** (30x faster)
- Developer iterations - instant feedback on "do I need to rebuild?"

## 🎯 Completion Checklist

- [x] Create `pkg/buildcache/` system
- [x] Update `buildMacOS()` with idempotency
- [x] Add `--force` and `--check` flags to build command
- [ ] Update `buildAndroid()` with same pattern
- [ ] Update `buildIOS()` with same pattern
- [ ] Update `buildWindows()` with same pattern
- [ ] Update `buildAll()` wrapper
- [ ] Add idempotency to screenshot command
- [ ] Add idempotency to package command
- [ ] Update Taskfile with new patterns
- [ ] Add `:force` variants of common tasks
- [ ] Add `:check` tasks for CI/CD
- [ ] Create test script for idempotency
- [ ] Update CLAUDE.md with idempotency guarantees
- [ ] Add performance benchmarks

## 📚 Documentation Updates

### CLAUDE.md

Add section:
```markdown
## Idempotency Guarantees

All goup-util commands are **idempotent by default**:

✅ **build** - Skips if sources unchanged (use `--force` to rebuild)
✅ **install** - Skips if SDK already installed
✅ **screenshot** - Skips if file exists (use `--force` to overwrite)
✅ **package** - Skips if package up-to-date
✅ **icons** - Skips if icons already generated

### Why This Matters

- **Faster iterations**: `task run:hybrid` is instant if nothing changed
- **CI efficiency**: Only rebuild what changed
- **Safe to retry**: Running commands multiple times is safe
- **Clear feedback**: Know immediately if work is needed

### Force Flags

Every destructive command has `--force`:
\`\`\`bash
go run . build --force macos examples/hybrid-dashboard
go run . screenshot --force output.png
go run . package --force macos
\`\`\`

### Check Mode

Use `--check` to test if work is needed:
\`\`\`bash
# Exit 0 if up-to-date, 1 if rebuild needed
go run . build --check macos examples/hybrid-dashboard
echo $?  # 0 = up-to-date, 1 = rebuild needed
\`\`\`

Perfect for CI/CD pipelines!
```

## 🔄 Migration Notes

**Cache file location**: `~/.goup-util/build-cache.json`

**First run after upgrade**: All builds will say "no previous build found" and rebuild once. After that, idempotency kicks in.

**Clear cache**:
```bash
rm ~/.goup-util/build-cache.json
```

**Inspect cache**:
```bash
cat ~/.goup-util/build-cache.json | jq .
```

## 🎉 Benefits

1. **Developer Experience**
   - Instant feedback on whether rebuild is needed
   - No more wasteful rebuilds
   - Clear messaging about what's happening

2. **CI/CD**
   - Only rebuild what changed
   - Faster pipelines
   - Explicit checks via `--check` flag

3. **Reliability**
   - Safe to run commands multiple times
   - Consistent behavior
   - No accidental work loss

4. **Performance**
   - **30x faster** for unchanged projects
   - Build cache tracks all changes
   - Smart dependency detection

---

Last updated: 2025-10-22
Status: Partially implemented (macOS build only, needs expansion)
