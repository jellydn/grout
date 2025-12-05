package ui

import (
	"fmt"
	"grout/models"
	"grout/state"
	"path/filepath"
	"slices"
	"strings"
	"time"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
	"qlova.tech/sum"
)

type GameList struct {
	Platform     models.Platform
	Games        models.Items
	SearchFilter string
}

func InitGamesList(platform models.Platform, games models.Items, searchFilter string) GameList {
	var g models.Items

	if len(games) > 0 {
		g = games
	} else {
		process, err := gaba.ProcessMessage(fmt.Sprintf("Loading %s...", platform.Name), gaba.ProcessMessageOptions{ShowThemeBackground: true}, func() (interface{}, error) {
			var err error
			g, err = loadGamesList(platform)
			return g, err
		})
		if err != nil {
			return GameList{}
		}

		g = process.(models.Items)
	}

	state.SetCurrentFullGamesList(g)

	return GameList{
		Platform:     platform,
		Games:        g,
		SearchFilter: searchFilter,
	}
}

func (gl GameList) Name() sum.Int[models.ScreenName] {
	return models.ScreenNames.GameList
}

func (gl GameList) Draw() (game interface{}, exitCode int, e error) {
	title := gl.Platform.Name

	itemList := gl.Games

	for idx, _ := range itemList {
		if itemList[idx].DisplayName == "" {
			itemList[idx].DisplayName = strings.ReplaceAll(itemList[idx].Filename, filepath.Ext(itemList[idx].Filename), "")
		}
	}

	slices.SortFunc(itemList, func(a, b models.Item) int {
		return strings.Compare(strings.ToLower(a.DisplayName), strings.ToLower(b.DisplayName))
	})

	if gl.SearchFilter != "" {
		title = fmt.Sprintf("[Search: \"%s\"] | %s", gl.SearchFilter, gl.Platform.Name)
		itemList = filterList(itemList, gl.SearchFilter)
	}

	if len(itemList) == 0 {
		if gl.SearchFilter != "" {
			gaba.ProcessMessage(
				fmt.Sprintf("No results found for \"%s\"", gl.SearchFilter),
				gaba.ProcessMessageOptions{ShowThemeBackground: true},
				func() (interface{}, error) {
					time.Sleep(time.Second * 2)
					return nil, nil
				},
			)
		} else {
			gaba.ProcessMessage(
				fmt.Sprintf("No games found for %s", gl.Platform.Name),
				gaba.ProcessMessageOptions{ShowThemeBackground: true},
				func() (interface{}, error) {
					time.Sleep(time.Second * 2)
					return nil, nil
				},
			)
		}
		return nil, 404, nil
	}

	var itemEntries []gaba.MenuItem
	for _, game := range itemList {
		itemEntries = append(itemEntries, gaba.MenuItem{
			Text:     game.DisplayName,
			Selected: false,
			Focused:  false,
			Metadata: game,
		})
	}

	options := gaba.DefaultListOptions(title, itemEntries)
	options.SmallTitle = true
	options.EnableAction = true
	options.EnableMultiSelect = true
	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: "Back"},
		{ButtonName: "X", HelpText: "Search"},
		{ButtonName: "Select", HelpText: "Multi"},
		{ButtonName: "A", HelpText: "Select"},
	}
	options.SelectedIndex = state.GetAppState().LastSelectedIndex
	options.VisibleStartIndex = max(0, state.GetAppState().LastSelectedIndex-state.GetAppState().LastSelectedPosition)

	res, err := gaba.List(options)
	if err != nil {
		return nil, 2, err
	}

	if res.Action == gaba.ListActionSelected {

		var selections models.Items
		for _, idx := range res.Selected {
			selections = append(selections, res.Items[idx].Metadata.(models.Item))
		}

		state.SetLastSelectedPosition(res.Selected[0], res.VisiblePosition)

		return selections, 0, nil
	} else if res.Action == gaba.ListActionTriggered {
		return nil, 4, nil
	}

	return nil, 2, err
}

func loadGamesList(platform models.Platform) (games models.Items, e error) {
	logger := gaba.GetLogger()

	items, err := FetchListStateless(platform)
	if err != nil {
		logger.Error("Error downloading Item List", "error", err)
	}

	if len(items) == 0 {
		return nil, nil
	}

	slices.SortFunc(items, func(a, b models.Item) int {
		return strings.Compare(strings.ToLower(a.Filename), strings.ToLower(b.Filename))
	})

	return items, nil

}
