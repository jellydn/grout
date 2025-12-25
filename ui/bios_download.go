package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"grout/constants"
	"grout/romm"
	"grout/utils"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
	"github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/i18n"
)

type BIOSDownloadInput struct {
	Config   utils.Config
	Host     romm.Host
	Platform romm.Platform
}

type BIOSDownloadOutput struct {
	Platform romm.Platform
}

type BIOSDownloadScreen struct{}

func NewBIOSDownloadScreen() *BIOSDownloadScreen {
	return &BIOSDownloadScreen{}
}

func (s *BIOSDownloadScreen) Execute(config utils.Config, host romm.Host, platform romm.Platform) BIOSDownloadOutput {
	result, err := s.draw(BIOSDownloadInput{
		Config:   config,
		Host:     host,
		Platform: platform,
	})

	if err != nil {
		gaba.GetLogger().Error("BIOS download failed", "error", err)
		return BIOSDownloadOutput{Platform: platform}
	}

	return result.Value
}

func (s *BIOSDownloadScreen) draw(input BIOSDownloadInput) (ScreenResult[BIOSDownloadOutput], error) {
	logger := gaba.GetLogger()

	output := BIOSDownloadOutput{
		Platform: input.Platform,
	}

	biosFiles := utils.GetBIOSFilesForPlatform(input.Platform.Slug)

	if len(biosFiles) == 0 {
		logger.Info("No BIOS files required for platform", "platform", input.Platform.Name)
		// TODO: Show user-friendly message
		return back(output), nil
	}

	var menuItems []gaba.MenuItem
	var biosStatusMap = make(map[string]utils.BIOSFileStatus)

	for _, biosFile := range biosFiles {
		status := utils.CheckBIOSFileStatus(biosFile, input.Platform.Slug)
		biosStatusMap[biosFile.FileName] = status

		var statusIndicator string
		var statusText string
		switch status.Status {
		case utils.BIOSStatusValid:
			statusIndicator = "✓"
			statusText = i18n.GetString("bios_status_ready")
		case utils.BIOSStatusInvalidHash:
			statusIndicator = "⚠"
			statusText = i18n.GetString("bios_status_wrong_version")
		case utils.BIOSStatusNoHashToVerify:
			statusIndicator = "?"
			statusText = i18n.GetString("bios_status_unverified")
		case utils.BIOSStatusMissing:
			statusIndicator = "✗"
			statusText = i18n.GetString("bios_status_not_installed")
		}

		optionalText := ""
		if biosFile.Optional {
			optionalText = " (Optional)"
		}

		displayText := fmt.Sprintf("%s %s%s - %s", statusIndicator, biosFile.FileName, optionalText, statusText)

		menuItems = append(menuItems, gaba.MenuItem{
			Text:     displayText,
			Selected: status.Status == utils.BIOSStatusMissing || status.Status == utils.BIOSStatusInvalidHash,
			Focused:  false,
			Metadata: biosFile,
		})
	}

	options := gaba.DefaultListOptions(fmt.Sprintf("%s - BIOS Files", input.Platform.Name), menuItems)
	options.StartInMultiSelectMode = true
	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: "Back"},
		{ButtonName: "A", HelpText: "Toggle Selection"},
		{ButtonName: "Start", HelpText: "Download Selected"},
	}

	sel, err := gaba.List(options)
	if err != nil {
		logger.Error("BIOS selection failed", "error", err)
		return back(output), err
	}

	if sel.Action != gaba.ListActionSelected || len(sel.Selected) == 0 {
		return back(output), nil
	}

	var selectedBIOSFiles []constants.BIOSFile
	for _, idx := range sel.Selected {
		biosFile := sel.Items[idx].Metadata.(constants.BIOSFile)
		selectedBIOSFiles = append(selectedBIOSFiles, biosFile)
	}

	logger.Debug("Selected BIOS files for download", "count", len(selectedBIOSFiles))
	for i, biosFile := range selectedBIOSFiles {
		logger.Debug("Selected BIOS file",
			"index", i,
			"filename", biosFile.FileName,
			"relativePath", biosFile.RelativePath,
			"md5", biosFile.MD5Hash,
			"optional", biosFile.Optional)
	}

	client := utils.GetRommClient(input.Host, input.Config.ApiTimeout)
	firmwareList, err := client.GetFirmware(input.Platform.ID)
	if err != nil {
		logger.Error("Failed to fetch firmware from RomM", "error", err, "platform_id", input.Platform.ID)
		// TODO: Show user-friendly error message
		return back(output), nil
	}

	logger.Debug("Fetched firmware from RomM", "count", len(firmwareList), "platform_id", input.Platform.ID)

	for i, fw := range firmwareList {
		logger.Debug("RomM firmware entry",
			"index", i,
			"id", fw.ID,
			"filename", fw.FileName,
			"filepath", fw.FilePath,
			"fullpath", fw.FullPath,
			"size", fw.FileSizeBytes,
			"md5", fw.MD5Hash,
			"verified", fw.IsVerified,
			"has_download_url", fw.DownloadURL != "",
			"download_url", fw.DownloadURL)
	}

	downloads, locationToBIOSMap := s.buildDownloads(input.Host, selectedBIOSFiles, firmwareList)

	if len(downloads) == 0 {
		logger.Warn("No BIOS files available in RomM")
		// TODO: Show user-friendly message
		return back(output), nil
	}

	headers := make(map[string]string)
	headers["Authorization"] = input.Host.BasicAuthHeader()

	res, err := gaba.DownloadManager(downloads, headers, gaba.DownloadManagerOptions{
		AutoContinue: false,
	})
	if err != nil {
		logger.Error("BIOS download failed", "error", err)
		return back(output), err
	}

	logger.Debug("Download results", "completed", len(res.Completed), "failed", len(res.Failed))

	successCount := 0
	warningCount := 0
	for _, download := range res.Completed {
		biosFile := locationToBIOSMap[download.Location]

		data, err := os.ReadFile(download.Location)
		if err != nil {
			logger.Error("Failed to read downloaded BIOS file", "file", biosFile.FileName, "error", err)
			continue
		}

		if biosFile.MD5Hash != "" {
			isValid, actualHash := utils.VerifyBIOSFileMD5(data, biosFile.MD5Hash)
			if !isValid {
				logger.Warn("MD5 hash mismatch for BIOS file",
					"file", biosFile.FileName,
					"expected", biosFile.MD5Hash,
					"actual", actualHash)
				warningCount++
			}
		}

		if err := utils.SaveBIOSFile(biosFile, input.Platform.Slug, data); err != nil {
			logger.Error("Failed to save BIOS file", "file", biosFile.FileName, "error", err)
			continue
		}

		os.Remove(download.Location)
		successCount++
		logger.Debug("Successfully saved BIOS file", "file", biosFile.FileName, "path", utils.GetBIOSFilePath(biosFile, input.Platform.Slug))
	}

	if successCount > 0 && warningCount == 0 {
		logger.Info("BIOS download complete", "success", successCount)
	} else if successCount > 0 && warningCount > 0 {
		logger.Warn("BIOS download complete with warnings",
			"success", successCount,
			"warnings", warningCount)
	} else if len(res.Failed) > 0 {
		logger.Error("BIOS download failed", "failed", len(res.Failed))
	}

	// TODO: Show user-friendly completion message

	return back(output), nil
}

func (s *BIOSDownloadScreen) buildDownloads(host romm.Host, biosFiles []constants.BIOSFile, firmwareList []romm.Firmware) ([]gaba.Download, map[string]constants.BIOSFile) {
	var downloads []gaba.Download
	locationToBIOSMap := make(map[string]constants.BIOSFile)

	logger := gaba.GetLogger()
	baseURL := host.URL()

	firmwareByFileName := make(map[string]romm.Firmware)
	firmwareByFilePath := make(map[string]romm.Firmware)
	firmwareByBaseName := make(map[string]romm.Firmware)

	for _, fw := range firmwareList {
		firmwareByFileName[fw.FileName] = fw
		firmwareByFilePath[fw.FilePath] = fw
		baseName := filepath.Base(fw.FilePath)
		firmwareByBaseName[baseName] = fw
	}

	for _, biosFile := range biosFiles {
		var firmware romm.Firmware
		var found bool
		var matchStrategy string

		firmware, found = firmwareByFileName[biosFile.FileName]
		if found {
			matchStrategy = "exact_filename"
		}
		if !found {
			firmware, found = firmwareByFilePath[biosFile.RelativePath]
			if found {
				matchStrategy = "relative_path"
			}
		}
		if !found {
			firmware, found = firmwareByBaseName[biosFile.FileName]
			if found {
				matchStrategy = "basename_filename"
			}
		}
		if !found {
			baseName := filepath.Base(biosFile.RelativePath)
			firmware, found = firmwareByBaseName[baseName]
			if found {
				matchStrategy = "basename_relativepath"
			}
		}

		if !found {
			logger.Warn("BIOS file not found in RomM firmware list",
				"file", biosFile.FileName,
				"relativePath", biosFile.RelativePath)
			continue
		}

		logger.Debug("Matched BIOS file with RomM firmware",
			"biosFile", biosFile.FileName,
			"firmware", firmware.FileName,
			"strategy", matchStrategy)

		downloadURL := baseURL + firmware.DownloadURL

		tempPath := filepath.Join(utils.TempDir(), fmt.Sprintf("bios_%s", biosFile.FileName))

		downloads = append(downloads, gaba.Download{
			URL:         downloadURL,
			Location:    tempPath,
			DisplayName: biosFile.FileName,
		})

		locationToBIOSMap[tempPath] = biosFile

		logger.Debug("Added BIOS file to download queue",
			"file", biosFile.FileName,
			"url", downloadURL,
			"size", firmware.FileSizeBytes)
	}

	return downloads, locationToBIOSMap
}
