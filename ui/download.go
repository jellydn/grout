package ui

import (
	"encoding/base64"
	"grout/models"
	"grout/utils"
	"net/url"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"grout/romm"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

type DownloadInput struct {
	Config        models.Config
	Host          models.Host
	Platform      romm.Platform
	SelectedGames []romm.Rom
	AllGames      []romm.Rom
	SearchFilter  string
}

type DownloadOutput struct {
	DownloadedGames []romm.Rom
	Platform        romm.Platform
	AllGames        []romm.Rom
	SearchFilter    string
}

type DownloadScreen struct{}

func NewDownloadScreen() *DownloadScreen {
	return &DownloadScreen{}
}

func (s *DownloadScreen) Execute(config models.Config, host models.Host, platform romm.Platform, selectedGames []romm.Rom, allGames []romm.Rom, searchFilter string) DownloadOutput {
	result, err := s.Draw(DownloadInput{
		Config:        config,
		Host:          host,
		Platform:      platform,
		SelectedGames: selectedGames,
		AllGames:      allGames,
		SearchFilter:  searchFilter,
	})

	if err != nil {
		gaba.GetLogger().Error("Download failed", "error", err)
		return DownloadOutput{
			AllGames:     allGames,
			Platform:     platform,
			SearchFilter: searchFilter,
		}
	}

	if result.ExitCode == gaba.ExitCodeSuccess && len(result.Value.DownloadedGames) > 0 {
		gaba.GetLogger().Debug("Successfully downloaded games", "count", len(result.Value.DownloadedGames))
	}

	return result.Value
}

func (s *DownloadScreen) Draw(input DownloadInput) (ScreenResult[DownloadOutput], error) {
	logger := gaba.GetLogger()

	output := DownloadOutput{
		Platform:     input.Platform,
		AllGames:     input.AllGames,
		SearchFilter: input.SearchFilter,
	}

	downloads := s.buildDownloads(input.Config, input.Host, input.Platform, input.SelectedGames)

	headers := make(map[string]string)
	auth := input.Host.Username + ":" + input.Host.Password
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	headers["Authorization"] = authHeader

	logger.Debug("RomM Auth Header", "header", authHeader)

	slices.SortFunc(downloads, func(a, b gaba.Download) int {
		return strings.Compare(strings.ToLower(a.DisplayName), strings.ToLower(b.DisplayName))
	})

	logger.Debug("Starting ROM download", "downloads", downloads)

	res, err := gaba.DownloadManager(downloads, headers, input.Config.DownloadArt)
	if err != nil {
		logger.Error("Error downloading", "error", err)
		return WithCode(output, gaba.ExitCodeError), err
	}

	if len(res.Failed) > 0 {
		for _, g := range downloads {
			failedMatch := slices.ContainsFunc(res.Failed, func(de gaba.DownloadError) bool {
				return de.Download.DisplayName == g.DisplayName
			})
			if failedMatch {
				utils.DeleteFile(g.Location)
			}
		}
	}

	if len(res.Completed) == 0 {
		return WithCode(output, gaba.ExitCodeError), nil
	}

	downloadedGames := make([]romm.Rom, 0, len(res.Completed))
	for _, g := range input.SelectedGames {
		if slices.ContainsFunc(res.Completed, func(d gaba.Download) bool {
			return d.DisplayName == g.Name
		}) {
			downloadedGames = append(downloadedGames, g)
		}
	}

	output.DownloadedGames = downloadedGames
	return Success(output), nil
}

func (s *DownloadScreen) buildDownloads(config models.Config, host models.Host, platform romm.Platform, games []romm.Rom) []gaba.Download {
	downloads := make([]gaba.Download, 0, len(games))

	for _, g := range games {
		// For collections, use each game's platform info; for platforms, use the passed platform
		gamePlatform := platform
		if platform.ID == 0 && g.PlatformID != 0 {
			// Construct platform from game's platform info (happens when viewing collections)
			gamePlatform = romm.Platform{
				ID:   g.PlatformID,
				Slug: g.PlatformSlug,
				Name: g.PlatformDisplayName,
			}
		}

		romDirectory := utils.GetPlatformRomDirectory(config, gamePlatform)
		downloadLocation := ""

		sourceURL := ""

		if g.Multi {
			// TODO Fill this shit out
		} else {
			downloadLocation = filepath.Join(romDirectory, g.Files[0].FileName)
			sourceURL, _ = url.JoinPath(host.URL(), "/api/roms/", strconv.Itoa(g.ID), "content", g.Files[0].FileName)
		}

		downloads = append(downloads, gaba.Download{
			URL:         sourceURL,
			Location:    downloadLocation,
			DisplayName: g.Name,
			Timeout:     config.DownloadTimeout,
		})
	}

	return downloads
}
