package resources

import (
	"sam-api/models"
)

//Models for logical model resources envelopes
type (
	// request
	DictionaryAccountBscsRequestResource struct {
		Data models.DictionaryAccountBscs `json:"data"`
	}

	// reply with feedback, one objects just created
	DictionaryAccountBscsReplyResource struct {
		Data models.DictionaryAccountBscs `json:"data"`
	}

	// reply with many objects
	DictionaryAccountBscssReplyResource struct {
		Count int64                          `json:"count"`
		Data  []models.DictionaryAccountBscs `json:"data"`
	}
)
