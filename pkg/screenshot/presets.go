// Package screenshot provides App Store screenshot presets
package screenshot

// Preset defines a screenshot size preset for app stores
type Preset struct {
	Name        string
	Width       int
	Height      int
	Description string
	Store       string // "ios", "macos", "android", "windows"
}

// Presets contains all available App Store screenshot presets
var Presets = map[string]Preset{
	// iOS App Store
	"iphone-6.9": {
		Name:        "iPhone 6.9\" (iPhone 16 Pro Max)",
		Width:       1320,
		Height:      2868,
		Description: "iPhone 16 Pro Max portrait",
		Store:       "ios",
	},
	"iphone-6.7": {
		Name:        "iPhone 6.7\" (iPhone 15 Plus)",
		Width:       1290,
		Height:      2796,
		Description: "iPhone 15 Plus portrait",
		Store:       "ios",
	},
	"iphone-6.5": {
		Name:        "iPhone 6.5\" (iPhone 15 Pro Max)",
		Width:       1284,
		Height:      2778,
		Description: "iPhone 15 Pro Max portrait",
		Store:       "ios",
	},
	"iphone-5.5": {
		Name:        "iPhone 5.5\" (iPhone 8 Plus)",
		Width:       1242,
		Height:      2208,
		Description: "iPhone 8 Plus portrait",
		Store:       "ios",
	},
	"ipad-12.9": {
		Name:        "iPad Pro 12.9\" (6th gen)",
		Width:       2048,
		Height:      2732,
		Description: "iPad Pro 12.9\" portrait",
		Store:       "ios",
	},
	"ipad-11": {
		Name:        "iPad Pro 11\" (4th gen)",
		Width:       1668,
		Height:      2388,
		Description: "iPad Pro 11\" portrait",
		Store:       "ios",
	},

	// macOS App Store
	"macos-retina": {
		Name:        "macOS Retina",
		Width:       2880,
		Height:      1800,
		Description: "macOS Retina display (recommended)",
		Store:       "macos",
	},
	"macos-retina-2k": {
		Name:        "macOS Retina 2K",
		Width:       2560,
		Height:      1600,
		Description: "macOS Retina 2K display",
		Store:       "macos",
	},
	"macos-standard": {
		Name:        "macOS Standard",
		Width:       1280,
		Height:      800,
		Description: "macOS standard display (minimum)",
		Store:       "macos",
	},

	// Android Play Store
	"android-phone": {
		Name:        "Android Phone",
		Width:       1080,
		Height:      1920,
		Description: "Android phone portrait (1080p)",
		Store:       "android",
	},
	"android-phone-landscape": {
		Name:        "Android Phone Landscape",
		Width:       1920,
		Height:      1080,
		Description: "Android phone landscape (1080p)",
		Store:       "android",
	},
	"android-tablet-7": {
		Name:        "Android Tablet 7\"",
		Width:       1200,
		Height:      1920,
		Description: "Android 7\" tablet portrait",
		Store:       "android",
	},
	"android-tablet-10": {
		Name:        "Android Tablet 10\"",
		Width:       1600,
		Height:      2560,
		Description: "Android 10\" tablet portrait",
		Store:       "android",
	},

	// Windows Store
	"windows-hd": {
		Name:        "Windows HD",
		Width:       1920,
		Height:      1080,
		Description: "Windows Full HD (recommended)",
		Store:       "windows",
	},
	"windows-min": {
		Name:        "Windows Minimum",
		Width:       1366,
		Height:      768,
		Description: "Windows minimum size",
		Store:       "windows",
	},
	"windows-xbox": {
		Name:        "Windows Xbox",
		Width:       1920,
		Height:      1080,
		Description: "Xbox console",
		Store:       "windows",
	},
}

// GetPreset returns a preset by name
func GetPreset(name string) (Preset, bool) {
	preset, ok := Presets[name]
	return preset, ok
}

// ListPresets returns all presets, optionally filtered by store
func ListPresets(store string) []Preset {
	var result []Preset
	for _, p := range Presets {
		if store == "" || p.Store == store {
			result = append(result, p)
		}
	}
	return result
}

// ListPresetNames returns all preset names, optionally filtered by store
func ListPresetNames(store string) []string {
	var result []string
	for name, p := range Presets {
		if store == "" || p.Store == store {
			result = append(result, name)
		}
	}
	return result
}
