# Code Signing Status

## Current Status: ⚠️ Unsigned Binaries

**All goup-util binaries are currently UNSIGNED.**

This means users will see security warnings when downloading and running the tool.

### What This Means

#### macOS
- **First run**: macOS Gatekeeper will block the binary
- **Warning**: "goup-util cannot be opened because it is from an unidentified developer"
- **Workaround**: Our bootstrap script automatically runs `xattr -d com.apple.quarantine`
- **Manual bypass**: System Settings → Privacy & Security → "Allow Anyway"

#### Windows
- **First run**: Windows SmartScreen will warn about unsigned publisher
- **Warning**: "Windows protected your PC"
- **Workaround**: Our bootstrap script runs `Unblock-File`
- **Manual bypass**: Click "More info" → "Run anyway"

#### Linux
- No signing required for executables

---

## Bootstrap Scripts Handle This

Our bootstrap scripts **automatically handle unsigned binaries**:

### macOS (`macos-bootstrap.sh`)
```bash
# Removes quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/goup-util

# Alternative commands if issues persist:
sudo spctl --add /usr/local/bin/goup-util
```

### Windows (`windows-bootstrap.ps1`)
```powershell
# Unblocks the downloaded file
Unblock-File -Path $env:USERPROFILE\goup-util.exe

# Adds to PATH for easy access
```

---

## Future: Code Signing

### Why Sign?

**Benefits of signed binaries**:
1. ✅ No security warnings for users
2. ✅ Better trust and credibility
3. ✅ Required for macOS Gatekeeper
4. ✅ Required for Windows SmartScreen
5. ✅ Enables automatic updates without warnings

### What's Needed

#### macOS Signing

**Option 1: Apple Developer ID** (Recommended)
- **Cost**: $99/year
- **Process**:
  1. Join Apple Developer Program
  2. Request Developer ID Application certificate
  3. Sign binaries: `codesign --sign "Developer ID Application: Name" goup-util`
  4. Notarize: `xcrun notarytool submit goup-util.zip`
- **Result**: No Gatekeeper warnings, automatic updates work

**Option 2: Ad-hoc Signing** (Current alternative)
- **Cost**: Free
- **Process**: `codesign --sign - goup-util`
- **Limitation**: Gatekeeper still blocks, users must manually allow

#### Windows Signing

**Option 1: EV Code Signing Certificate** (Recommended)
- **Cost**: ~$300-500/year
- **Vendors**: DigiCert, Sectigo, GlobalSign
- **Process**:
  1. Purchase EV Code Signing Certificate
  2. Sign binary: `signtool sign /f cert.pfx /p password goup-util.exe`
- **Result**: No SmartScreen warnings, builds reputation

**Option 2: Standard Code Signing**
- **Cost**: ~$100-200/year
- **Process**: Same as EV
- **Limitation**: SmartScreen warnings until reputation builds

#### Linux
- No signing required for executables
- Package repositories (apt, yum) handle verification

---

## Implementation Timeline

### Phase 1: Current (No Signing)
- ✅ Bootstrap scripts handle unsigned binaries
- ✅ Clear documentation and warnings
- ✅ Manual workarounds provided

### Phase 2: macOS Signing (Q1 2026)
**When**: After securing funding or sponsorship
**Steps**:
1. Purchase Apple Developer Program membership
2. Generate Developer ID Application certificate
3. Update GitHub Actions to sign binaries
4. Notarize macOS binaries
5. Update bootstrap script to expect signed binaries

**CI/CD Integration**:
```yaml
# .github/workflows/release.yml
- name: Sign macOS binary
  run: |
    codesign --sign "${{ secrets.APPLE_SIGNING_IDENTITY }}" \
      --entitlements pkg/packaging/templates/macos-entitlements.plist.tmpl \
      --options runtime \
      goup-util-darwin-arm64

- name: Notarize
  run: |
    xcrun notarytool submit goup-util-darwin-arm64.zip \
      --apple-id "${{ secrets.APPLE_ID }}" \
      --password "${{ secrets.APP_SPECIFIC_PASSWORD }}" \
      --team-id "${{ secrets.TEAM_ID }}"
```

### Phase 3: Windows Signing (Q2 2026)
**When**: After securing code signing certificate
**Steps**:
1. Purchase EV Code Signing Certificate
2. Store certificate in GitHub Secrets
3. Update GitHub Actions to sign binaries
4. Update bootstrap script

**CI/CD Integration**:
```yaml
# .github/workflows/release.yml
- name: Sign Windows binary
  run: |
    signtool sign /f cert.pfx /p "${{ secrets.WINDOWS_CERT_PASSWORD }}" \
      /tr http://timestamp.digicert.com \
      /td sha256 \
      /fd sha256 \
      goup-util-windows-amd64.exe
```

---

## Cost Summary

| Item | Cost | Frequency | Priority |
|------|------|-----------|----------|
| Apple Developer Program | $99 | Annual | High |
| Windows EV Certificate | $300-500 | Annual | Medium |
| Windows Standard Certificate | $100-200 | Annual | Low |
| **Total Annual** | **$400-700** | Annual | - |

**Alternative**: Community sponsorship or donations to cover signing costs.

---

## User Impact

### Current (Unsigned)

**First-time users**:
1. Run bootstrap script (handles security warnings)
2. Or manually bypass security warnings
3. Binary works normally after initial bypass

**Updates via `goup-util self upgrade`**:
1. Security warnings reappear on each update
2. Must manually bypass again
3. Annoying but workable

### After Signing

**First-time users**:
1. Download runs without warnings
2. Binary executes immediately
3. No manual intervention needed

**Updates**:
1. `goup-util self upgrade` works seamlessly
2. No security warnings
3. Automatic background updates possible

---

## Temporary Workarounds

### For Developers

**macOS**:
```bash
# Remove quarantine from all goup-util binaries
find /usr/local/bin -name "goup-util*" -exec xattr -d com.apple.quarantine {} \;

# Disable Gatekeeper system-wide (NOT recommended)
sudo spctl --master-disable
```

**Windows**:
```powershell
# Unblock all downloaded binaries
Get-ChildItem -Path $env:USERPROFILE -Filter "goup-util*" | Unblock-File

# Disable SmartScreen (NOT recommended)
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "SmartScreenEnabled" -Value "Off"
```

### For CI/CD

**GitHub Actions**:
- Binaries built in GitHub Actions are not quarantined
- No signing needed for internal CI/CD testing
- Only affects downloaded binaries

**Internal Distribution**:
- Host binaries on internal server
- Add server to trusted sources
- Or sign with internal certificate

---

## Questions?

- **Why not self-signed certificates?**
  - Self-signed certificates are treated the same as unsigned
  - No benefit over current approach

- **Can users sign themselves?**
  - macOS: Users can sign with their own Developer ID
  - Windows: Users can sign with their own certificate
  - But this doesn't help other users

- **What about open source projects?**
  - Many open source projects are unsigned
  - This is a known issue in the ecosystem
  - Users are accustomed to workarounds

- **When will signing happen?**
  - When funding is available ($400-700/year)
  - Or when Windows support becomes critical
  - Or when enough users request it

---

## See Also

- `scripts/macos-bootstrap.sh` - Handles unsigned macOS binaries
- `scripts/windows-bootstrap.ps1` - Handles unsigned Windows binaries
- `cmd/self.go` - Self-upgrade functionality
- [Apple Code Signing](https://developer.apple.com/support/code-signing/)
- [Windows Code Signing](https://docs.microsoft.com/en-us/windows/win32/seccrypto/cryptography-tools)
