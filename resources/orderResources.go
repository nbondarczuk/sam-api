package resources

import (
	"sam-api/models"
)

//Models for logical model resources envelopes
type (
	// Request
	OrderRequestResource struct {
		Data models.Order `json:"data"`
	}

	// reply with feedback, one objects just created
	OrderReplyResource struct {
		Data models.Order `json:"data"`
	}

	// Reply
	OrdersReplyResource struct {
		Count int64          `json:"count"`
		Data  []models.Order `json:"data"`
	}

	// Reply with logs
	OrderLogsReplyResource struct {
		Count int64               `json:"count"`
		Data  []models.OrderLog `json:"data"`
	}
)
