package models

import "strings"

type Host struct {
	DisplayName string `yaml:"display_name,omitempty" json:"display_name,omitempty"`
	RootURI     string `yaml:"root_uri,omitempty" json:"root_uri,omitempty"`
	Port        int    `yaml:"port,omitempty" json:"port,omitempty"`

	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	Platforms Platforms `yaml:"-" json:"-"`
}
type Hosts []Host

func (h Host) ToLoggable() map[string]any {
	temp := map[string]any{
		"display_name": h.DisplayName,
		"root_uri":     h.RootURI,
		"port":         h.Port,
		"username":     h.Username,
		"password":     strings.Repeat("*", len(h.Password)),
		"platforms":    h.Platforms.ToLoggable(),
	}

	return temp
}
