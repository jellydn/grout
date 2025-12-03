package ui

import (
	"fmt"
	"grout/models"
	"grout/state"
	"grout/utils"
	"time"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"qlova.tech/sum"
)

type SettingsScreen struct {
}

func InitSettingsScreen() SettingsScreen {
	return SettingsScreen{}
}

func (s SettingsScreen) Name() sum.Int[models.ScreenName] {
	return models.ScreenNames.Settings
}

func (s SettingsScreen) Draw() (settings interface{}, exitCode int, e error) {
	logger := gabagool.GetLogger()

	appState := state.GetAppState()

	items := []gabagool.ItemWithOptions{
		{
			Item: gabagool.MenuItem{
				Text: "Download Art",
			},
			Options: []gabagool.Option{
				{DisplayName: "True", Value: true},
				{DisplayName: "False", Value: false},
			},
			SelectedOption: func() int {
				if appState.Config.DownloadArt {
					return 0
				}
				return 1
			}(),
		},

		{
			Item: gabagool.MenuItem{
				Text: "API Timeout",
			},
			Options: []gabagool.Option{
				{DisplayName: "15 Seconds", Value: "15s"},
				{DisplayName: "30 Seconds", Value: "30s"},
				{DisplayName: "45 Seconds", Value: "45s"},
				{DisplayName: "60 Seconds", Value: "60s"},
				{DisplayName: "75 Seconds", Value: "75s"},
				{DisplayName: "90 Seconds", Value: "90s"},
				{DisplayName: "120 Seconds", Value: "120s"},
				{DisplayName: "180 Seconds", Value: "180s"},
				{DisplayName: "240 Seconds", Value: "240s"},
				{DisplayName: "300 Seconds", Value: "300s"},
			},
			SelectedOption: func() int {
				seconds := int(appState.Config.ApiTimeout.Seconds())
				switch seconds {
				case 15:
					return 0
				case 30:
					return 1
				case 45:
					return 2
				case 60:
					return 3
				case 75:
					return 4
				case 90:
					return 5
				case 120:
					return 6
				case 180:
					return 7
				case 240:
					return 8
				case 300:
					return 9
				default:
					return 0
				}
			}(),
		},
		{
			Item: gabagool.MenuItem{
				Text: "Download Timeout",
			},
			Options: []gabagool.Option{
				{DisplayName: "15 Minutes", Value: "15m"},
				{DisplayName: "30 Minutes", Value: "30m"},
				{DisplayName: "45 Minutes", Value: "45m"},
				{DisplayName: "60 Minutes", Value: "60m"},
				{DisplayName: "75 Minutes", Value: "75m"},
				{DisplayName: "90 Minutes", Value: "90m"},
				{DisplayName: "105 Minutes", Value: "105m"},
				{DisplayName: "120 Minutes", Value: "120m"},
			},
			SelectedOption: func() int {
				minutes := int(appState.Config.DownloadTimeout.Minutes())
				switch minutes {
				case 15:
					return 0
				case 30:
					return 1
				case 45:
					return 2
				case 60:
					return 3
				case 75:
					return 4
				case 90:
					return 5
				case 105:
					return 6
				case 120:
					return 7
				default:
					return 0
				}
			}(),
		},

		{
			Item: gabagool.MenuItem{
				Text: "Use Title As Filename",
			},
			Options: []gabagool.Option{
				{DisplayName: "True", Value: true},
				{DisplayName: "False", Value: false},
			},
			SelectedOption: func() int {
				if appState.Config.UseTitleAsFilename {
					return 0
				}
				return 1
			}(),
		},

		{
			Item: gabagool.MenuItem{
				Text: "Unzip Downloads",
			},
			Options: []gabagool.Option{
				{DisplayName: "True", Value: true},
				{DisplayName: "False", Value: false},
			},
			SelectedOption: func() int {
				if appState.Config.UnzipDownloads {
					return 0
				}
				return 1
			}(),
		},
		{
			Item: gabagool.MenuItem{
				Text: "Group BIN / CUE",
			},
			Options: []gabagool.Option{
				{DisplayName: "True", Value: true},
				{DisplayName: "False", Value: false},
			},
			SelectedOption: func() int {
				if appState.Config.GroupBinCue {
					return 0
				}
				return 1
			}(),
		},
		{
			Item: gabagool.MenuItem{
				Text: "Group Multi-Disc",
			},
			Options: []gabagool.Option{
				{DisplayName: "True", Value: true},
				{DisplayName: "False", Value: false},
			},
			SelectedOption: func() int {
				if appState.Config.GroupMultiDisc {
					return 0
				}
				return 1
			}(),
		},
		{
			Item: gabagool.MenuItem{
				Text: "Log Level",
			},
			Options: []gabagool.Option{
				{DisplayName: "Debug", Value: "DEBUG"},
				{DisplayName: "Error", Value: "ERROR"},
			},
			SelectedOption: func() int {
				switch appState.Config.LogLevel {
				case "DEBUG":
					return 0
				case "ERROR":
					return 1
				}
				return 0
			}(),
		},
	}

	result, err := gabagool.OptionsList(
		"Grout Settings",
		gabagool.OptionListSettings{FooterHelpItems: []gabagool.FooterHelpItem{
			{ButtonName: "B", HelpText: "Cancel"},
			{ButtonName: "←→", HelpText: "Cycle"},
			{ButtonName: "Start", HelpText: "Save"},
		}},
		items,
	)

	if err != nil {
		fmt.Println("Error showing options list:", err)
		return
	}

	if result.IsSome() {
		newSettingOptions := result.Unwrap().Items

		for _, option := range newSettingOptions {
			switch option.Item.Text {
			case "Download Art":
				appState.Config.DownloadArt = option.SelectedOption == 0
			case "API Timeout":
				timeoutValues := []time.Duration{
					15 * time.Second,
					30 * time.Second,
					45 * time.Second,
					60 * time.Second,
					75 * time.Second,
					90 * time.Second,
					120 * time.Second,
					180 * time.Second,
					240 * time.Second,
					300 * time.Second,
				}
				appState.Config.ApiTimeout = timeoutValues[option.SelectedOption]
			case "Download Timeout":
				timeoutValues := []time.Duration{
					15 * time.Minute,
					30 * time.Minute,
					45 * time.Minute,
					60 * time.Minute,
					75 * time.Minute,
					90 * time.Minute,
					105 * time.Minute,
					120 * time.Minute,
				}
				appState.Config.DownloadTimeout = timeoutValues[option.SelectedOption]
			case "Use Title As Filename":
				appState.Config.UseTitleAsFilename = option.SelectedOption == 0
			case "Unzip Downloads":
				appState.Config.UnzipDownloads = option.SelectedOption == 0
			case "Group BIN / CUE":
				appState.Config.GroupBinCue = option.SelectedOption == 0
			case "Group Multi-Disc":
				appState.Config.GroupMultiDisc = option.SelectedOption == 0
			case "Log Level":
				logLevelValue := option.Options[option.SelectedOption].Value.(string)
				appState.Config.LogLevel = logLevelValue
			}
		}

		err := utils.SaveConfig(appState.Config)
		if err != nil {
			logger.Error("Error saving config", "error", err)
			return nil, 0, err
		}

		state.UpdateAppState(appState)

		return result, 0, nil
	}

	return nil, 2, nil
}
