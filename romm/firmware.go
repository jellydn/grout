package romm

import (
	"strconv"
	"time"
)

type Firmware struct {
	ID          int       `json:"id"`
	PlatformID  int       `json:"platform_id"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	FileHash    string    `json:"file_hash"`
	Description *string   `json:"description"`
	Version     *string   `json:"version"`
	Required    bool      `json:"required"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) GetFirmware(platformID *int) ([]Firmware, error) {
	params := map[string]string{}

	if platformID != nil {
		params["platform_id"] = strconv.Itoa(*platformID)
	}

	var firmware []Firmware
	path := "/api/firmware" + buildQueryString(params)
	err := c.doRequest("GET", path, nil, &firmware)
	return firmware, err
}
