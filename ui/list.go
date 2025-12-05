package ui

import (
	"grout/client"
	"grout/models"
	"grout/state"
	"slices"
	"strings"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

func FetchListStateless(platform models.Platform) (models.Items, error) {
	logger := gaba.GetLogger()
	config := state.GetAppState().Config

	logger.Debug("Fetching Item List",
		"host", platform.Host.ToLoggable())

	c := client.NewRomMClient(platform.Host, config.ApiTimeout)

	defer func(client client.RomMClient) {
		err := client.Close()
		if err != nil {
			logger.Error("Unable to close client", "error", err)
		}
	}(*c)

	items, err := c.ListDirectory(platform.ID)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.Item, 0, len(items))
	for _, item := range items {
		if !strings.HasPrefix(item.Filename, ".") {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func filterList(itemList []models.Item, filter string) []models.Item {
	var result []models.Item

	for _, item := range itemList {
		if strings.Contains(strings.ToLower(item.DisplayName), strings.ToLower(filter)) {
			result = append(result, item)
		}
	}

	slices.SortFunc(result, func(a, b models.Item) int {
		return strings.Compare(strings.ToLower(a.DisplayName), strings.ToLower(b.DisplayName))
	})

	return result
}
