package models

import (
	"time"
)

type (
	DictionaryAccountSap struct {
		Account       string    `json:"sapOfiAccount" db:"SAP_OFI_ACCOUNT,size:32,primarykey"`
		Name          string    `json:"name" db:"NAME,size:255"`
		Status        string    `json:"status" db:"STATUS,size:8"`
		EntryDate     time.Time `json:"-" db:"ENTRY_DATE"`
		EntryDateStr  string    `json:"entryDate,omitempty" db:"-"`
		EntryOwner    string    `json:"entryOwner,omitempty" db:"ENTRY_OWNER,size:16"`
		UpdateDate    time.Time `json:"-" db:"UPDATE_DATE"`
		UpdateDateStr string    `json:"updateDate,omitempty" db:"-"`
		UpdateOwner   string    `json:"-" db:"UPDATE_OWNER,size:16"`
		RecVersion    int       `json:"recVersion" db:"REC_VERSION"`		
	}
)
