package artwork

import (
	"grout/internal/fileutil"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

func GetCacheDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join(os.TempDir(), ".cache", "artwork")
	}
	return filepath.Join(wd, ".cache", "artwork")
}

func ClearCache() error {
	cacheDir := GetCacheDir()

	if !fileutil.FileExists(cacheDir) {
		return nil
	}

	return os.RemoveAll(cacheDir)
}

func HasCache() bool {
	cacheDir := GetCacheDir()

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return false
	}

	return len(entries) > 0
}

func GetCachePath(platformFSSlug string, romID int) string {
	return filepath.Join(GetCacheDir(), platformFSSlug, strconv.Itoa(romID)+".png")
}

func Exists(platformFSSlug string, romID int) bool {
	return fileutil.FileExists(GetCachePath(platformFSSlug, romID))
}

func EnsureCacheDir(platformFSSlug string) error {
	dir := filepath.Join(GetCacheDir(), platformFSSlug)
	return os.MkdirAll(dir, 0755)
}

func ValidateCache() {
	go func() {
		logger := gaba.GetLogger()
		cacheDir := GetCacheDir()

		platformDirs, err := os.ReadDir(cacheDir)
		if err != nil {
			return
		}

		removed := 0
		for _, platformDir := range platformDirs {
			if !platformDir.IsDir() {
				continue
			}

			platformPath := filepath.Join(cacheDir, platformDir.Name())
			files, err := os.ReadDir(platformPath)
			if err != nil {
				continue
			}

			for _, file := range files {
				if file.IsDir() || filepath.Ext(file.Name()) != ".png" {
					continue
				}

				path := filepath.Join(platformPath, file.Name())
				f, err := os.Open(path)
				if err != nil {
					continue
				}

				_, err = png.DecodeConfig(f)
				f.Close()
				if err != nil {
					os.Remove(path)
					removed++
				}
			}
		}

		if removed > 0 {
			logger.Debug("Removed corrupted artwork files", "count", removed)
		}
	}()
}
