package models

import (
	"time"
)

type (
	Account struct {
		Status           string    `json:"status,omitempty" db:"STATUS,size:1,primarykey"`
		ReleaseId        string    `json:"releaseId,omitempty" db:"RELEASE_ID,size:8,primarykey"`
		BscsAccount      string    `json:"bscsAccount" db:"BSCS_ACCOUNT,size:32,primarykey"`
		OfiSapAccount    string    `json:"ofiSapAccount,omitempty" db:"OFI_SAP_ACCOUNT,size:32"`
		ValidFromDateStr string    `json:"validFromDate,omitempty" db:"-"`
		ValidFromDate    time.Time `json:"-" db:"VALID_FROM_DATE"`
		VatCodeInd       string    `json:"vatCodeInd,omitempty" db:"VAT_CODE_IND,size:1"`
		OfiSapWbsCode    string    `json:"ofiSapWbsCode,omitempty" db:"OFI_SAP_WBS_CODE,size:32"`
		CitMarkerVatFlag int       `json:"citMarkerVatFlag,omitempty" db:"CIT_MARKER_VAT_FLAG"`
		EntryDate        time.Time `json:"-" db:"ENTRY_DATE"`
		EntryDateStr     string    `json:"entryDate,omitempty" db:"-"`
		EntryOwner       string    `json:"entryOwner,omitempty" db:"ENTRY_OWNER,size:16"`
		UpdateDate       time.Time `json:"-" db:"UPDATE_DATE"`
		UpdateDateStr    string    `json:"updateDate,omitempty" db:"-"`
		UpdateOwner      string    `json:"-" db:"UPDATE_OWNER,size:16"`
		ReleaseDate      time.Time `json:"-" db:"RELEASE_DATE"`
		ReleaseDateStr   string    `json:"releaseDate,omitempty" db:"-"`
		ReleaseOwner     string    `json:"-" db:"RELEASE_OWNER,size:16"`
		RecVersion       int       `json:"recVersion" db:"REC_VERSION"`
	}

	Accounts struct {
		Data []Account
	}

	AccountLog struct {
		OpCode string    `json:"opcode" db:"OPCODE, size:1"`
		OpDate time.Time `json:"opdate" db:"OPDATE"`
		Account
	}
)
