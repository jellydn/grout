package ui

import (
	"grout/models"

	"qlova.tech/sum"
)

type PlatformMappingScreen struct {
}

func InitPlatformMappingScreen() PlatformMappingScreen {
	return PlatformMappingScreen{}
}

func (p PlatformMappingScreen) Name() sum.Int[models.ScreenName] {
	return models.ScreenNames.SettingsPlatformMapping
}

func (p PlatformMappingScreen) Draw() (settings interface{}, exitCode int, e error) {
	//logger := gabagool.GetLogger()
	//
	//appState := state.GetAppState()

	return nil, 0, nil
}
