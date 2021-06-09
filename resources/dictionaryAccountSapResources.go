package resources

import (
	"sam-api/models"
)

//Models for logical model resources envelopes
type (
	// request
	DictionaryAccountSapRequestResource struct {
		Data models.DictionaryAccountSap `json:"data"`
	}

	// reply with feedback, one objects just created
	DictionaryAccountSapReplyResource struct {
		Data models.DictionaryAccountSap `json:"data"`
	}

	// reply with many objects
	DictionaryAccountSapsReplyResource struct {
		Count int64                         `json:"count"`
		Data  []models.DictionaryAccountSap `json:"data"`
	}
)
