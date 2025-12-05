package models

type AppState struct {
	Config      *Config
	HostIndices map[string]int

	CurrentFullGamesList Items
	LastSelectedIndex    int
	LastSelectedPosition int
}
