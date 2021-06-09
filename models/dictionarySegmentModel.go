package models

import (
	"time"
)

type (
	DictionarySegment struct {
		CsTradeRef    string    `json:"csTradeRef" db:"CSTRADEREF,size:8,primarykey"`
		SegmCategory  string    `json:"segmCategory" db:"SEGM_CATEGORY,size:4"`
		EntryDate     time.Time `json:"-" db:"ENTRY_DATE"`
		EntryDateStr  string    `json:"entryDate,omitempty" db:"-"`
		EntryOwner    string    `json:"entryOwner,omitempty" db:"ENTRY_OWNER,size:16"`
		UpdateDate    time.Time `json:"-" db:"UPDATE_DATE"`
		UpdateDateStr string    `json:"updateDate,omitempty" db:"-"`
		UpdateOwner   string    `json:"-" db:"UPDATE_OWNER,size:16"`
		RecVersion    int       `json:"recVersion" db:"REC_VERSION"`		
	}
)
