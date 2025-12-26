package ui

import (
	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

type LanguageSelectionScreen struct{}

func NewLanguageSelectionScreen() *LanguageSelectionScreen {
	return &LanguageSelectionScreen{}
}

func (s *LanguageSelectionScreen) Draw() (string, error) {
	options := []gaba.SelectionOption{
		{DisplayName: "English", Value: "en"},
		{DisplayName: "Español", Value: "es"},
		{DisplayName: "Français", Value: "fr"},
	}

	result, err := gaba.SelectionMessage(
		"Language / Idioma / Langue",
		options,
		[]gaba.FooterHelpItem{
			{ButtonName: "←→", HelpText: "Select"},
			{ButtonName: "A", HelpText: "Confirm"},
		},
		gaba.SelectionMessageSettings{},
	)

	if err != nil {
		return "", err
	}

	return result.SelectedValue.(string), nil
}
