package ui

import (
	"errors"
	"grout/models"
	"grout/romm"
	"grout/utils"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

type CollectionSelectionInput struct {
	Host                 models.Host
	LastSelectedIndex    int
	LastSelectedPosition int
}

type CollectionSelectionOutput struct {
	SelectedCollection   romm.Collection
	LastSelectedIndex    int
	LastSelectedPosition int
}

type CollectionSelectionScreen struct{}

func NewCollectionSelectionScreen() *CollectionSelectionScreen {
	return &CollectionSelectionScreen{}
}

func (s *CollectionSelectionScreen) Draw(input CollectionSelectionInput) (ScreenResult[CollectionSelectionOutput], error) {
	output := CollectionSelectionOutput{
		LastSelectedIndex:    input.LastSelectedIndex,
		LastSelectedPosition: input.LastSelectedPosition,
	}

	// Fetch collections from RomM
	rc := utils.GetRommClient(input.Host)
	collections, err := rc.GetCollections()
	if err != nil {
		return WithCode(output, gaba.ExitCodeError), err
	}

	if len(collections) == 0 {
		return WithCode(output, gaba.ExitCode(404)), nil
	}

	var menuItems []gaba.MenuItem
	for _, collection := range collections {
		menuItems = append(menuItems, gaba.MenuItem{
			Text:     collection.Name,
			Selected: false,
			Focused:  false,
			Metadata: collection,
		})
	}

	footerItems := []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: "Back"},
		{ButtonName: "A", HelpText: "Select"},
	}

	options := gaba.DefaultListOptions("Collections", menuItems)
	options.FooterHelpItems = footerItems
	options.SelectedIndex = input.LastSelectedIndex
	options.VisibleStartIndex = max(0, input.LastSelectedIndex-input.LastSelectedPosition)

	sel, err := gaba.List(options)
	if err != nil {
		if errors.Is(err, gaba.ErrCancelled) {
			return Back(output), nil
		}
		return WithCode(output, gaba.ExitCodeError), err
	}

	switch sel.Action {
	case gaba.ListActionSelected:
		collection := sel.Items[sel.Selected[0]].Metadata.(romm.Collection)

		output.SelectedCollection = collection
		output.LastSelectedIndex = sel.Selected[0]
		output.LastSelectedPosition = sel.VisiblePosition
		return Success(output), nil
	}

	return WithCode(output, gaba.ExitCodeBack), nil
}
