# Screen Recording Permissions on macOS

**Quick Reference for Users**

## Why You Need to Grant Permission

goup-util uses the screenshot command for:
- 📸 Capturing app screenshots for documentation
- 🧪 Testing visual output
- 📊 Creating README examples

macOS 10.15+ (Catalina and later) requires explicit permission for any app that captures screen content.

## Grant Permission (One-Time Setup)

### Method 1: Let the App Trigger It

1. Run any screenshot command:
   ```bash
   CGO_ENABLED=1 go run . screenshot test.png
   ```

2. macOS will show a dialog:
   > **"Terminal" would like to record this computer's screen.**
   >
   > [Deny] [Allow]

3. Click **Allow**

4. Done! Permission is saved.

### Method 2: Manual Setup

If you missed the dialog or denied it:

1. Open **System Settings** (or System Preferences on older macOS)

2. Go to **Privacy & Security**

3. Click **Screen Recording** in the sidebar

4. Check the box next to your terminal app:
   - ☑️ Terminal.app
   - ☑️ iTerm2
   - ☑️ VSCode
   - ☑️ Cursor
   - ☑️ (whatever you're using)

5. **Restart your terminal app**

6. Try the screenshot command again

## For Different Terminals

### Terminal.app (Default macOS)
✅ Add "Terminal" to Screen Recording list

### iTerm2
✅ Add "iTerm" to Screen Recording list

### VS Code Integrated Terminal
✅ Add "Visual Studio Code" or "Code" to Screen Recording list

### Cursor
✅ Add "Cursor" to Screen Recording list

### Warp
✅ Add "Warp" to Screen Recording list

## Using the Packaged App

If you build a packaged version:

```bash
# Build the .app bundle
task package:macos

# Open it
open .dist/goup-util.app
```

Then grant permission to **goup-util** (not your terminal).

## Verifying Permission

### Check if permission is granted:

```bash
# Try to capture display info
CGO_ENABLED=1 go run . screenshot --info
```

**If permission is granted**, you'll see:
```
Found 1 display(s):

Display 0:
  Position: (0, 0)
  Size: 1728x1117
```

**If permission is NOT granted**, you'll see:
```
Error: screenshot failed: Capture image not found.

Note: On macOS 10.15+, grant Screen Recording permission:
System Settings → Privacy & Security → Screen Recording
```

## Common Issues

### "Permission denied" after granting

**Solution**: Restart your terminal app completely (Cmd+Q, then reopen)

### Permission keeps resetting

**Cause**: You're rebuilding the app with different code signatures

**Solution**: Use the packaged version:
```bash
task package:macos
```

This creates a properly signed app with a consistent bundle ID, so macOS remembers the permission.

### "Terminal is not in the list"

**Solution**: First run a screenshot command to make it appear, then grant permission.

## For Distribution

If you're distributing goup-util to other users:

1. **Package it properly**:
   ```bash
   task package:macos:dmg
   ```

2. **Include in your README**:
   ```markdown
   ## macOS Setup

   After installation, you need to grant Screen Recording permission:

   1. Run goup-util once (it will trigger the permission dialog)
   2. Or manually: System Settings → Privacy & Security → Screen Recording
   3. Check the box next to goup-util
   ```

3. **Users only need to do this ONCE** (unless they revoke permission)

## Why No Automatic Permission?

Unlike some other permissions (Camera, Microphone), macOS does **NOT** allow apps to automatically get screen recording access, even with entitlements.

This is a security feature. **All apps must**:
1. Be properly code-signed
2. Request permission (which triggers the dialog)
3. Wait for user to manually grant access

There's no way around this for third-party apps!

## Summary

✅ **One-time setup**: Grant Screen Recording permission
✅ **Remember**: Permission is saved per app (or terminal)
✅ **Packaged apps**: Bundle ID must stay consistent
✅ **Clear errors**: goup-util tells you exactly what to do

---

See also:
- [Screenshot Integration](SCREENSHOT.md) - Full screenshot features
- [macOS Packaging](PACKAGING-MACOS.md) - Code signing and distribution
