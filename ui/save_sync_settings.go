package ui

import (
	"errors"
	"grout/cfw"
	"grout/internal"
	"grout/utils"
	"sort"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
	"github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/i18n"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type SaveSyncSettingsInput struct {
	Config *internal.Config
	CFW    cfw.CFW
}

type SaveSyncSettingsOutput struct {
	Config *internal.Config
}

type SaveSyncSettingsScreen struct{}

func NewSaveSyncSettingsScreen() *SaveSyncSettingsScreen {
	return &SaveSyncSettingsScreen{}
}

func (s *SaveSyncSettingsScreen) Draw(input SaveSyncSettingsInput) (ScreenResult[SaveSyncSettingsOutput], error) {
	config := input.Config
	output := SaveSyncSettingsOutput{Config: config}

	items := s.buildMenuItems(config)

	if len(items) == 0 {
		gaba.GetLogger().Warn("No platforms configured for save sync settings")
		return back(output), nil
	}

	result, err := gaba.OptionsList(
		i18n.Localize(&goi18n.Message{ID: "save_sync_settings_title", Other: "Save Sync Settings"}, nil),
		gaba.OptionListSettings{
			FooterHelpItems: []gaba.FooterHelpItem{
				FooterBack(),
				FooterCycle(),
				FooterSave(),
			},
			InitialSelectedIndex: 0,
			StatusBar:            utils.StatusBar(),
			SmallTitle:           true,
		},
		items,
	)

	if err != nil {
		if errors.Is(err, gaba.ErrCancelled) {
			return back(output), nil
		}
		gaba.GetLogger().Error("Save sync settings error", "error", err)
		return withCode(output, gaba.ExitCodeError), err
	}

	s.applySettings(config, result.Items)

	err = internal.SaveConfig(config)
	if err != nil {
		gaba.GetLogger().Error("Error saving save sync settings", "error", err)
		return withCode(output, gaba.ExitCodeError), err
	}

	return success(output), nil
}

func (s *SaveSyncSettingsScreen) buildMenuItems(config *internal.Config) []gaba.ItemWithOptions {
	items := make([]gaba.ItemWithOptions, 0)

	// Get all platform fsSlugs from directory mappings
	fsSlugs := make([]string, 0, len(config.DirectoryMappings))
	for fsSlug := range config.DirectoryMappings {
		fsSlugs = append(fsSlugs, fsSlug)
	}
	sort.Strings(fsSlugs)

	for _, fsSlug := range fsSlugs {
		saveDirectories := cfw.EmulatorFoldersForFSSlug(fsSlug)
		if len(saveDirectories) == 0 {
			continue
		}

		options := make([]gaba.Option, 0, len(saveDirectories)+1)

		// Add "Default" option first
		options = append(options, gaba.Option{
			DisplayName: i18n.Localize(&goi18n.Message{ID: "common_default", Other: "Default"}, nil),
			Value:       "",
		})

		// Add each emulator directory as an option
		for _, dir := range saveDirectories {
			options = append(options, gaba.Option{
				DisplayName: dir,
				Value:       dir,
			})
		}

		// Determine currently selected option
		selectedIndex := 0
		if config.SaveDirectoryMappings != nil {
			if currentMapping, ok := config.SaveDirectoryMappings[fsSlug]; ok && currentMapping != "" {
				for i, opt := range options {
					if val, ok := opt.Value.(string); ok && val == currentMapping {
						selectedIndex = i
						break
					}
				}
			}
		}

		items = append(items, gaba.ItemWithOptions{
			Item:           gaba.MenuItem{Text: fsSlug},
			Options:        options,
			SelectedOption: selectedIndex,
		})
	}

	return items
}

func (s *SaveSyncSettingsScreen) applySettings(config *internal.Config, items []gaba.ItemWithOptions) {
	if config.SaveDirectoryMappings == nil {
		config.SaveDirectoryMappings = make(map[string]string)
	}

	for _, item := range items {
		fsSlug := item.Item.Text
		if val, ok := item.Options[item.SelectedOption].Value.(string); ok {
			if val == "" {
				// Remove from map if set to default
				delete(config.SaveDirectoryMappings, fsSlug)
			} else {
				config.SaveDirectoryMappings[fsSlug] = val
			}
		}
	}
}
