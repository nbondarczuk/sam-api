package resources

import (
	"sam-api/models"
)

//Models for logical model resources envelopes
type (
	// Request
	AccountRequestResource struct {
		Data models.Account `json:"data"`
	}

	// reply with feedback, one objects just created
	AccountReplyResource struct {
		Data models.Account `json:"data"`
	}

	// Reply
	AccountsReplyResource struct {
		Count int64            `json:"count"`
		Data  []models.Account `json:"data"`
	}

	// Reply with logs
	AccountLogsReplyResource struct {
		Count int64               `json:"count"`
		Data  []models.AccountLog `json:"data"`
	}
)
