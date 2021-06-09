package models

import (
	"time"
)

type (
	DictionaryAccountBscs struct {
		Account       string    `json:"account" db:"GLACODE,size:30,primarykey"`
		Name          string    `json:"name" db:"GLADESC,size:60"`
		Type          string    `json:"type" db:"GLATYPE,size:1"`
		Active        string    `json:"active" db:"GLACTIVE,size:1"`
		EntryDate     time.Time `json:"-" db:"ENTRY_DATE"`
		EntryDateStr  string    `json:"entryDate,omitempty" db:"-"`
		EntryOwner    string    `json:"entryOwner,omitempty" db:"ENTRY_OWNER,size:16"`
		UpdateDate    time.Time `json:"-" db:"UPDATE_DATE"`
		UpdateDateStr string    `json:"updateDate,omitempty" db:"-"`
		UpdateOwner   string    `json:"-" db:"UPDATE_OWNER,size:16"`
	}
)
