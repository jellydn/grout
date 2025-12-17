package romm

import (
	"time"
)

type Firmware struct {
	ID          int       `json:"id"`
	PlatformID  int       `json:"platform_id"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	FileHash    string    `json:"file_hash"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Required    bool      `json:"required"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FirmwareOptions struct {
	PlatformID int `qs:"platform_id"`
}

func (fo FirmwareOptions) Valid() bool {
	return fo.PlatformID != 0
}

func (c *Client) getFirmware(platformID int) ([]Firmware, error) {
	var firmware []Firmware
	err := c.doRequest("GET", endpointFirmware, FirmwareOptions{PlatformID: platformID}, nil, &firmware)
	return firmware, err
}
