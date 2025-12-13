package utils

import (
	"fmt"
	"grout/constants"
	"os"
	"path/filepath"
	"strings"
	"time"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

const BACKUP_TIMESTAMP_FORMAT = "2006-01-02 15-04-05"

type LocalSave struct {
	Slug         string
	Path         string
	LastModified time.Time
}

func (lc LocalSave) TimestampedFilename() string {
	ext := filepath.Ext(lc.Path)
	base := strings.ReplaceAll(filepath.Base(lc.Path), ext, "")

	lm := lc.LastModified.Format(BACKUP_TIMESTAMP_FORMAT)

	return fmt.Sprintf("%s [%s]%s", base, lm, ext)
}

func (lc LocalSave) Backup() error {
	dest := filepath.Join(filepath.Dir(lc.Path), ".backup", lc.TimestampedFilename())
	return CopyFile(lc.Path, dest)
}

func GetSaveDirectoryForSlug(slug string, emulator string) (string, error) {
	logger := gaba.GetLogger()
	bsd := GetSaveDirectory()

	config, err := LoadConfig()
	if err == nil {
		if mapping, ok := config.DirectoryMappings[slug]; ok && mapping.SaveDirectory != "" {
			saveDir := filepath.Join(bsd, mapping.SaveDirectory)
			logger.Debug("Using config save directory", "slug", slug, "directory", mapping.SaveDirectory)

			if err := os.MkdirAll(saveDir, 0755); err != nil {
				logger.Error("Failed to create save directory", "path", saveDir, "error", err)
				return "", fmt.Errorf("failed to create save directory: %w", err)
			}

			return saveDir, nil
		}
	}

	var saveFolders []string

	switch GetCFW() {
	case constants.MuOS:
		saveFolders = constants.MuOSSaveDirectories[slug]
	case constants.NextUI:
		saveFolders = constants.NextUISaves[slug]
	}

	if len(saveFolders) == 0 {
		return "", fmt.Errorf("no save folder mapping for slug: %s", slug)
	}

	selectedFolder := saveFolders[0]
	if emulator != "" {
		for _, folder := range saveFolders {
			if strings.Contains(strings.ToLower(folder), strings.ToLower(emulator)) {
				selectedFolder = folder
				logger.Debug("Matched emulator to save folder", "emulator", emulator, "folder", folder)
				break
			}
		}
	}

	saveDir := filepath.Join(bsd, selectedFolder)

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		logger.Error("Failed to create save directory", "path", saveDir, "error", err)
		return "", fmt.Errorf("failed to create save directory: %w", err)
	}

	return saveDir, nil
}

func FindSaveFiles(slug string) []LocalSave {
	logger := gaba.GetLogger()

	bsd := GetSaveDirectory()
	var saveFolders []string

	switch GetCFW() {
	case constants.MuOS:
		saveFolders = constants.MuOSSaveDirectories[slug]
	case constants.NextUI:
		saveFolders = constants.NextUISaves[slug]
	}

	if len(saveFolders) == 0 {
		logger.Debug("No save folder mapping for slug", "slug", slug)
		return []LocalSave{}
	}

	var allSaveFiles []LocalSave

	for _, saveFolder := range saveFolders {
		sd := filepath.Join(bsd, saveFolder)

		if _, err := os.Stat(sd); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(sd)
		if err != nil {
			logger.Error("Failed to read save directory", "path", sd, "error", err)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				savePath := filepath.Join(sd, entry.Name())

				fileInfo, err := entry.Info()
				if err != nil {
					logger.Warn("Failed to get file info", "file", entry.Name(), "error", err)
					continue
				}

				saveFile := LocalSave{
					Slug:         slug,
					Path:         savePath,
					LastModified: fileInfo.ModTime(),
				}

				allSaveFiles = append(allSaveFiles, saveFile)
			}
		}

		logger.Debug("Found save files in directory", "path", sd, "count", len(entries))
	}

	logger.Debug("Found total save files", "slug", slug, "count", len(allSaveFiles))
	return allSaveFiles
}
