package resources

import (
	"sam-api/models"
)

//Models for logical model resources envelopes
type (
	// request
	DictionarySegmentRequestResource struct {
		Data models.DictionarySegment `json:"data"`
	}

	// reply with feedback, one objects just created
	DictionarySegmentReplyResource struct {
		Data models.DictionarySegment `json:"data"`
	}

	// reply with many objects
	DictionarySegmentsReplyResource struct {
		Count int64                      `json:"count"`
		Data  []models.DictionarySegment `json:"data"`
	}
)
