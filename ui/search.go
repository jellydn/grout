package ui

import (
	"errors"
	"grout/models"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
	"qlova.tech/sum"
)

type Search struct {
	Platform    models.Platform
	InitialText string
}

func InitSearch(platform models.Platform, initialText string) Search {
	return Search{
		Platform:    platform,
		InitialText: initialText,
	}
}

func (s Search) Name() sum.Int[models.ScreenName] {
	return models.ScreenNames.SearchBox
}

func (s Search) Draw() (value interface{}, exitCode int, e error) {
	res, err := gaba.Keyboard(s.InitialText)
	if err != nil {
		if !errors.Is(err, gaba.ErrCancelled) {
			gaba.GetLogger().Error("Error with keyboard", "error", err)
		}
		return nil, 2, err
	}

	return res.Text, 0, nil
}
