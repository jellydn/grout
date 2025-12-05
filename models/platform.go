package models

type Platform struct {
	ID   string `yaml:"-" json:"-"`
	Name string `yaml:"-" json:"-"`
	Slug string `yaml:"-" json:"-"`

	Host Host `json:"-"`
}

type Platforms []Platform

func (p Platform) ToLoggable() map[string]any {
	return map[string]any{
		"name":        p.Name,
		"platform_id": p.ID,
		"slug":        p.Slug,
	}
}

func (ps Platforms) ToLoggable() []map[string]any {
	var temp []map[string]any

	for _, p := range ps {
		temp = append(temp, p.ToLoggable())
	}

	return temp
}
