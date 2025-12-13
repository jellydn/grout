package ui

import gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"

type ScreenResult[T any] struct {
	Value    T
	ExitCode gaba.ExitCode
}

func Success[T any](value T) ScreenResult[T] {
	return ScreenResult[T]{
		Value:    value,
		ExitCode: gaba.ExitCodeSuccess,
	}
}

func Back[T any](value T) ScreenResult[T] {
	return ScreenResult[T]{
		Value:    value,
		ExitCode: gaba.ExitCodeBack,
	}
}

func WithCode[T any](value T, code gaba.ExitCode) ScreenResult[T] {
	return ScreenResult[T]{
		Value:    value,
		ExitCode: code,
	}
}
