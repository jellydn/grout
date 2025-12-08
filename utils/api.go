package utils

import (
	"grout/models"
	"grout/romm"
)

func GetRommClient(host models.Host) *romm.Client {
	return romm.NewClient(host.URL(), romm.WithBasicAuth(host.Username, host.Password))
}
