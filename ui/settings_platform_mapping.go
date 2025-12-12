package ui

import (
	"errors"
	"fmt"
	"grout/constants"
	"grout/utils"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"grout/romm"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

// PlatformMappingInput contains data needed to render the platform mapping screen
type PlatformMappingInput struct {
	Host           romm.Host
	ApiTimeout     time.Duration
	CFW            constants.CFW
	RomDirectory   string // Base ROM directory path
	AutoSelect     bool   // Auto-select "Create" option when no match found
	HideBackButton bool   // Hide the back/cancel button
}

// PlatformMappingOutput contains the result of the platform mapping screen
type PlatformMappingOutput struct {
	Mappings map[string]utils.DirectoryMapping
}

// PlatformMappingScreen displays platform to directory mapping configuration
type PlatformMappingScreen struct{}

func NewPlatformMappingScreen() *PlatformMappingScreen {
	return &PlatformMappingScreen{}
}

func (s *PlatformMappingScreen) Draw(input PlatformMappingInput) (ScreenResult[PlatformMappingOutput], error) {
	logger := gaba.GetLogger()
	output := PlatformMappingOutput{Mappings: make(map[string]utils.DirectoryMapping)}

	// Fetch RomM platforms
	rommPlatforms, err := s.fetchPlatforms(input)
	if err != nil {
		logger.Error("Error fetching RomM Platforms", "error", err)
		return WithCode(output, gaba.ExitCodeError), err
	}

	// Get local ROM directories
	romDirectories, err := s.getRomDirectories(input.RomDirectory)
	if err != nil {
		logger.Error("Error fetching ROM directories", "error", err)
		return WithCode(output, gaba.ExitCodeBack), err
	}

	// Build mapping options
	mappingOptions := s.buildMappingOptions(rommPlatforms, romDirectories, input)

	// Configure footer
	footerItems := []gaba.FooterHelpItem{
		{ButtonName: "←→", HelpText: "Cycle"},
		{ButtonName: "Start", HelpText: "Save"},
	}
	if !input.HideBackButton {
		footerItems = slices.Insert(footerItems, 0, gaba.FooterHelpItem{ButtonName: "B", HelpText: "Cancel"})
	}

	// Show options list
	result, err := gaba.OptionsList(
		"Rom Directory Mapping",
		gaba.OptionListSettings{
			FooterHelpItems:   footerItems,
			DisableBackButton: input.HideBackButton,
		},
		mappingOptions,
	)

	if err != nil {
		if errors.Is(err, gaba.ErrCancelled) {
			return Back(PlatformMappingOutput{}), nil
		}
		return WithCode(PlatformMappingOutput{}, gaba.ExitCodeError), err
	}

	// Build mappings from result
	output.Mappings = s.buildMappingsFromResult(result.Items)

	// Create directories for any "Create" options selected
	if err := s.createDirectories(output.Mappings, input.RomDirectory, romDirectories); err != nil {
		logger.Error("Error creating directories", "error", err)
		return WithCode(output, gaba.ExitCodeError), err
	}

	return Success(output), nil
}

func (s *PlatformMappingScreen) fetchPlatforms(input PlatformMappingInput) ([]romm.Platform, error) {
	client := romm.NewClient(
		input.Host.URL(),
		romm.WithBasicAuth(input.Host.Username, input.Host.Password),
		romm.WithTimeout(input.ApiTimeout),
	)
	return client.GetPlatforms()
}

func (s *PlatformMappingScreen) getRomDirectories(romDir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(romDir)
	if err != nil {
		gaba.ConfirmationMessage("ROM Directory Could Not Be Found!", []gaba.FooterHelpItem{
			{ButtonName: "B", HelpText: "Quit"},
		}, gaba.MessageOptions{})
		return nil, fmt.Errorf("failed to read ROM directory: %w", err)
	}

	var dirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry)
		}
	}

	return dirs, nil
}

func (s *PlatformMappingScreen) buildMappingOptions(
	platforms []romm.Platform,
	romDirectories []os.DirEntry,
	input PlatformMappingInput,
) []gaba.ItemWithOptions {
	options := make([]gaba.ItemWithOptions, 0, len(platforms))

	for _, platform := range platforms {
		platformOptions, selectedIndex := s.buildPlatformOptions(platform, romDirectories, input)

		options = append(options, gaba.ItemWithOptions{
			Item: gaba.MenuItem{
				Text:     platform.Name,
				Metadata: platform.Slug,
			},
			Options:        platformOptions,
			SelectedOption: selectedIndex,
		})
	}

	return options
}

func (s *PlatformMappingScreen) buildPlatformOptions(
	platform romm.Platform,
	romDirectories []os.DirEntry,
	input PlatformMappingInput,
) ([]gaba.Option, int) {
	// Start with "Skip" option
	options := []gaba.Option{{DisplayName: "Skip", Value: ""}}
	selectedIndex := 0

	// Get all valid CFW directories for this platform
	cfwDirectories := s.getCFWDirectoriesForPlatform(platform.Slug, input.CFW)

	// Add "Create" options for each valid CFW directory that doesn't exist
	createOptionAdded := false
	for _, cfwDir := range cfwDirectories {
		// Check if this directory already exists
		dirExists := false
		for _, romDir := range romDirectories {
			if s.directoriesMatch(cfwDir, romDir.Name(), input.CFW) {
				dirExists = true
				break
			}
		}

		if !dirExists {
			displayName := cfwDir
			if input.CFW == constants.NextUI {
				displayName = utils.ParseTag(cfwDir)
			}
			options = append(options, gaba.Option{
				DisplayName: fmt.Sprintf("Create '%s'", displayName),
				Value:       cfwDir,
			})
			createOptionAdded = true
		}
	}

	// Add existing ROM directories that match any valid CFW option for this platform
	for _, romDir := range romDirectories {
		dirName := romDir.Name()

		// Only add if this directory is valid for this platform
		if s.isValidDirectoryForPlatform(platform.Slug, dirName, input.CFW, cfwDirectories) {
			displayName := dirName
			if input.CFW == constants.NextUI {
				displayName = utils.ParseTag(dirName)
			}

			options = append(options, gaba.Option{
				DisplayName: fmt.Sprintf("/%s", displayName),
				Value:       dirName,
			})

			// Check if this directory matches the platform
			if s.directoryMatchesPlatform(platform, romDir.Name(), input.CFW) {
				selectedIndex = len(options) - 1
			}
		}
	}

	// Auto-select first "Create" option if no match found and auto-select enabled
	if selectedIndex == 0 && createOptionAdded && input.AutoSelect {
		selectedIndex = 1
	}

	return options, selectedIndex
}

func (s *PlatformMappingScreen) findMatchingDirectory(
	platform romm.Platform,
	romDirectories []os.DirEntry,
	cfw constants.CFW,
) int {
	for i, entry := range romDirectories {
		if s.directoryMatchesPlatform(platform, entry.Name(), cfw) {
			return i
		}
	}
	return -1
}

func (s *PlatformMappingScreen) directoryMatchesPlatform(
	platform romm.Platform,
	dirName string,
	cfw constants.CFW,
) bool {
	cfwSlug := utils.RomMSlugToCFW(platform.Slug)
	romFolderBase := utils.RomFolderBase(dirName)

	switch cfw {
	case constants.NextUI:
		return utils.ParseTag(cfwSlug) == romFolderBase
	default:
		return cfwSlug == romFolderBase
	}
}

func (s *PlatformMappingScreen) getCFWDirectoriesForPlatform(slug string, cfw constants.CFW) []string {
	switch cfw {
	case constants.MuOS:
		return constants.MuOSPlatforms[slug]
	case constants.NextUI:
		return constants.NextUIPlatforms[slug]
	default:
		return []string{}
	}
}

func (s *PlatformMappingScreen) getSaveDirectoriesForPlatform(slug string, cfw constants.CFW) []string {
	switch cfw {
	case constants.MuOS:
		return constants.MuOSSaveDirectories[slug]
	case constants.NextUI:
		return constants.NextUISaves[slug]
	default:
		return []string{}
	}
}

func (s *PlatformMappingScreen) directoriesMatch(dir1, dir2 string, cfw constants.CFW) bool {
	if cfw == constants.NextUI {
		return utils.ParseTag(dir1) == utils.ParseTag(dir2)
	}
	return dir1 == dir2
}

func (s *PlatformMappingScreen) isValidDirectoryForPlatform(slug, dirName string, cfw constants.CFW, cfwDirectories []string) bool {
	// Check if this directory matches any of the valid CFW directories for this platform
	for _, cfwDir := range cfwDirectories {
		if s.directoriesMatch(cfwDir, dirName, cfw) {
			return true
		}
	}
	return false
}

func (s *PlatformMappingScreen) getCreateDisplayName(slug string, cfw constants.CFW) string {
	displayName := utils.RomMSlugToCFW(slug)
	if cfw == constants.NextUI {
		displayName = utils.ParseTag(displayName)
	}
	return displayName
}

func (s *PlatformMappingScreen) buildMappingsFromResult(items []gaba.ItemWithOptions) map[string]utils.DirectoryMapping {
	mappings := make(map[string]utils.DirectoryMapping)
	cfw := utils.GetCFW()

	for _, item := range items {
		rommSlug := item.Item.Metadata.(string)
		relativePath := item.Options[item.SelectedOption].Value.(string)

		// Skip empty mappings
		if relativePath == "" {
			continue
		}

		// Auto-populate save directory based on the selected ROM directory
		saveDir := s.inferSaveDirectory(rommSlug, relativePath, cfw)

		mappings[rommSlug] = utils.DirectoryMapping{
			RomMSlug:      rommSlug,
			RelativePath:  relativePath,
			SaveDirectory: saveDir,
		}
	}

	return mappings
}

func (s *PlatformMappingScreen) inferSaveDirectory(slug, romDir string, cfw constants.CFW) string {
	saveDirectories := s.getSaveDirectoriesForPlatform(slug, cfw)

	// Try to match the ROM directory to a save directory
	for _, saveDir := range saveDirectories {
		if s.directoriesMatch(romDir, saveDir, cfw) {
			return saveDir
		}
	}

	// If no match, use the first available save directory
	if len(saveDirectories) > 0 {
		return saveDirectories[0]
	}

	// Fallback to ROM directory name
	return romDir
}

func (s *PlatformMappingScreen) createDirectories(
	mappings map[string]utils.DirectoryMapping,
	romDirectory string,
	existingDirs []os.DirEntry,
) error {
	logger := gaba.GetLogger()

	// Build a map of existing directory names for quick lookup
	existingDirMap := make(map[string]bool)
	for _, dir := range existingDirs {
		existingDirMap[dir.Name()] = true
	}

	// Check each mapping and create directories for new ones
	for _, mapping := range mappings {
		// Skip if the directory already exists
		if existingDirMap[mapping.RelativePath] {
			continue
		}

		// Create the directory
		fullPath := filepath.Join(romDirectory, mapping.RelativePath)
		logger.Debug("Creating new ROM directory", "path", fullPath)

		if err := os.MkdirAll(fullPath, 0755); err != nil {
			logger.Error("Failed to create directory", "path", fullPath, "error", err)
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}

		logger.Info("Created ROM directory", "path", fullPath)
	}

	return nil
}
