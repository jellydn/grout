package ui

import (
	"fmt"
	"grout/models"
	"grout/utils"
	"time"

	"github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
	"qlova.tech/sum"
)

type DownloadArtScreen struct {
	Platform     models.Platform
	Games        models.Items
	SearchFilter string
}

func InitDownloadArtScreen(platform models.Platform, games models.Items, searchFilter string) models.Screen {
	return DownloadArtScreen{
		Platform:     platform,
		Games:        games,
		SearchFilter: searchFilter,
	}
}

func (a DownloadArtScreen) Name() sum.Int[models.ScreenName] {
	return models.ScreenNames.DownloadArt
}

func (a DownloadArtScreen) Draw() (value interface{}, exitCode int, e error) {
	var artPaths []string

	gabagool.ProcessMessage("Downloading art...",
		gabagool.ProcessMessageOptions{ShowThemeBackground: true}, func() (interface{}, error) {
			for _, game := range a.Games {
				artPath := utils.FindArt(a.Platform, game)

				if artPath != "" {
					artPaths = append(artPaths, artPath)
				}
			}
			return nil, nil
		})

	if len(artPaths) == 0 {
		gabagool.ProcessMessage("No art downloaded!",
			gabagool.ProcessMessageOptions{ShowThemeBackground: true}, func() (interface{}, error) {
				time.Sleep(time.Millisecond * 1500)
				return nil, nil
			})

		return nil, 404, nil
	} else if len(a.Games) > 1 {
		gabagool.ProcessMessage(fmt.Sprintf("Art downloaded for %d/%d games!", len(artPaths), len(a.Games)),
			gabagool.ProcessMessageOptions{ShowThemeBackground: true}, func() (interface{}, error) {
				time.Sleep(time.Millisecond * 1500)
				return nil, nil
			})
	}

	for _, artPath := range artPaths {
		result, err := gabagool.ConfirmationMessage("Use This Art?",
			[]gabagool.FooterHelpItem{
				{ButtonName: "B", HelpText: "No"},
				{ButtonName: "A", HelpText: "Yes"},
			},
			gabagool.MessageOptions{
				ImagePath: artPath,
			})

		if err != nil || result.Confirmed {
			utils.DeleteFile(artPath)
		}
	}

	time.Sleep(time.Millisecond * 100)

	return nil, 0, nil
}
