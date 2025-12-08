package utils

import (
	"grout/models"
)

func ShowCollections(host models.Host) bool {
	rc := GetRommClient(host)
	col, err := rc.GetCollections()

	return err == nil && len(col) > 0
}
