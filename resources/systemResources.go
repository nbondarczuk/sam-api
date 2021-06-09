package resources

import (
	"sam-api/models"
)

type (
	VersionResource struct {
		Status string         `json:"status"`
		Data   models.Version `json:"data"`
	}

	StatResource struct {
		Status string      `json:"status"`
		Data   models.Stat `json:"data"`
	}
)
