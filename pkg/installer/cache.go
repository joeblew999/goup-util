package installer

import (
	"encoding/json"
	"os"
)

// CacheEntry represents a single entry in the artifact cache.
type CacheEntry struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	InstallPath string `json:"installPath"`
}

// Cache represents the artifact cache.
type Cache struct {
	Entries map[string]CacheEntry `json:"entries"`
	path    string
}

// NewCache creates a new cache instance and loads it from the given path.
func NewCache(path string) (*Cache, error) {
	cache := &Cache{
		Entries: make(map[string]CacheEntry),
		path:    path,
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cache, nil // Cache file doesn't exist yet, return a new empty cache
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return cache, nil // Cache file is empty, return a new empty cache
	}

	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Save writes the cache to disk.
func (c *Cache) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}

// IsCached checks if an SDK is already in the cache and if its checksum matches.
func (c *Cache) IsCached(sdk *SDK) bool {
	entry, ok := c.Entries[sdk.Name]
	if !ok {
		return false
	}

	return entry.Checksum == sdk.Checksum && entry.InstallPath == sdk.InstallPath
}

// Add adds a new SDK to the cache.
func (c *Cache) Add(sdk *SDK) {
	c.Entries[sdk.Name] = CacheEntry{
		Name:        sdk.Name,
		Version:     sdk.Version,
		Checksum:    sdk.Checksum,
		InstallPath: sdk.InstallPath,
	}
}
