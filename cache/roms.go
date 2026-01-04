package cache

import (
	"encoding/json"
	"grout/internal/fileutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

func stripExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

type RomCacheEntry struct {
	RomID    int       `json:"rom_id"`
	RomName  string    `json:"rom_name"`
	CachedAt time.Time `json:"cached_at"`
}

type PlatformRomCache struct {
	Entries map[string]RomCacheEntry `json:"entries"`
}

type RomCache struct {
	platforms map[string]*PlatformRomCache // keyed by platform fs_slug
	mu        sync.RWMutex
}

var (
	romCache     *RomCache
	romCacheOnce sync.Once
)

func GetRomCacheDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join(os.TempDir(), ".cache", "roms")
	}
	return filepath.Join(wd, ".cache", "roms")
}

func getRomCacheFilePath(fsSlug string) string {
	return filepath.Join(GetRomCacheDir(), fsSlug+".json")
}

func getRomCache() *RomCache {
	romCacheOnce.Do(func() {
		romCache = &RomCache{
			platforms: make(map[string]*PlatformRomCache),
		}
	})
	return romCache
}

func (rc *RomCache) getPlatformCache(fsSlug string) *PlatformRomCache {
	rc.mu.RLock()
	pc, exists := rc.platforms[fsSlug]
	rc.mu.RUnlock()

	if exists {
		return pc
	}

	pc = rc.loadPlatform(fsSlug)

	rc.mu.Lock()
	rc.platforms[fsSlug] = pc
	rc.mu.Unlock()

	return pc
}

func (rc *RomCache) loadPlatform(fsSlug string) *PlatformRomCache {
	logger := gaba.GetLogger()
	pc := &PlatformRomCache{
		Entries: make(map[string]RomCacheEntry),
	}

	data, err := os.ReadFile(getRomCacheFilePath(fsSlug))
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Debug("Failed to read ROM cache", "fsSlug", fsSlug, "error", err)
		}
		return pc
	}

	if err := json.Unmarshal(data, pc); err != nil {
		logger.Debug("Failed to parse ROM cache", "fsSlug", fsSlug, "error", err)
		return &PlatformRomCache{Entries: make(map[string]RomCacheEntry)}
	}

	if pc.Entries == nil {
		pc.Entries = make(map[string]RomCacheEntry)
	}

	logger.Debug("Loaded ROM cache", "fsSlug", fsSlug, "entries", len(pc.Entries))
	return pc
}

func (rc *RomCache) savePlatform(fsSlug string, pc *PlatformRomCache) error {
	logger := gaba.GetLogger()

	if err := os.MkdirAll(GetRomCacheDir(), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(getRomCacheFilePath(fsSlug), data, 0644); err != nil {
		return err
	}

	logger.Debug("Saved ROM cache", "fsSlug", fsSlug, "entries", len(pc.Entries))
	return nil
}

func GetCachedRomIDByFilename(fsSlug, filename string) (int, string, bool) {
	rc := getRomCache()
	pc := rc.getPlatformCache(fsSlug)

	key := stripExtension(filename)

	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if entry, ok := pc.Entries[key]; ok {
		return entry.RomID, entry.RomName, true
	}

	return 0, "", false
}

func StoreRomID(fsSlug, filename string, romID int, romName string) {
	logger := gaba.GetLogger()
	rc := getRomCache()
	pc := rc.getPlatformCache(fsSlug)

	key := stripExtension(filename)

	rc.mu.Lock()
	pc.Entries[key] = RomCacheEntry{
		RomID:    romID,
		RomName:  romName,
		CachedAt: time.Now(),
	}

	entriesCopy := make(map[string]RomCacheEntry, len(pc.Entries))
	for k, v := range pc.Entries {
		entriesCopy[k] = v
	}
	rc.mu.Unlock()

	pcCopy := &PlatformRomCache{Entries: entriesCopy}
	if err := rc.savePlatform(fsSlug, pcCopy); err != nil {
		logger.Debug("Failed to save ROM cache", "fsSlug", fsSlug, "error", err)
	}
}

func ClearRomCache() error {
	cacheDir := GetRomCacheDir()

	if !fileutil.FileExists(cacheDir) {
		return nil
	}

	rc := getRomCache()
	rc.mu.Lock()
	rc.platforms = make(map[string]*PlatformRomCache)
	rc.mu.Unlock()

	return os.RemoveAll(cacheDir)
}

func HasRomCache() bool {
	cacheDir := GetRomCacheDir()

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return false
	}

	return len(entries) > 0
}
