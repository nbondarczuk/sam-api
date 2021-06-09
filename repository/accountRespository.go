package repository

import (
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"
	_ "gopkg.in/goracle.v2"

	"sam-api/common"
	"sam-api/models"
)

//
// Pepository being handled by request
//
type AccountRepository struct {
	Repository
}

//
// Creates new repository using existing db connection
//
func NewAccountRepository(user string, trans bool) (r *AccountRepository, err error) {
	log.Printf("Creating new repository for user: %s", user)

	if db, err := common.GetDbSession(); err != nil {
		return nil, err
	} else {
		dbmap := initRepository(db)
		dbmap.AddTableWithName(models.Account{}, "SAP_ACCOUNTS").
			SetKeys(false, "STATUS").
			SetKeys(false, "RELEASE_ID").
			SetKeys(false, "BSCS_ACCOUNT")
		dbmap.AddTableWithName(models.AccountLog{}, "SAP_ACCOUNTS_LOG").
			SetKeys(false, "BSCS_ACCOUNT")
		r = &AccountRepository{
			Repository{
				Owner: user,
				Db:    db,
				Dbmap: dbmap,
			},
		}
		if trans {
			r.t, err = dbmap.Begin()
			if err != nil {
				return nil, err
			}
		}
		r.m.Lock()
	}

	return
}

func (r *AccountRepository) Close() {
	r.m.Unlock()
}

func (r *AccountRepository) Commit() {
	if r.t != nil {
		err := r.t.Commit()
		if err != nil {
			log.Printf("Commit error: %s", err.Error())
		}
	}
}

func (r *AccountRepository) Rollback() {
	if r.t != nil {
		err := r.t.Rollback()
		if err != nil {
			log.Printf("Rollback error: %s", err.Error())
		}
	}
}

//
// Insert new record to the resource SAP_ACCOUNTS
//
func (r *AccountRepository) Create(a *models.Account) (err error) {
	log.Printf("Inserting to SAP_ACCOUNTS: %s %#v", r.Owner, *a)

	// default value
	if a.ReleaseId == "" {
		a.ReleaseId = "0"
	}

	// default value
	if a.Status == "" {
		a.Status = "W"
	}

	if a.ValidFromDateStr != "" {
		var err error		
		a.ValidFromDate, err = time.Parse(common.CutOffDateFormat, a.ValidFromDateStr)
		if err != nil {
			return err
		}
	}

	a.EntryDate = time.Now()
	a.EntryOwner = r.Owner

	if r.t != nil {
		err = r.t.Insert(a)
	} else {
		err = r.Dbmap.Insert(a)
	}

	if err != nil {
		return fmt.Errorf("Error in insert to SAP_ACCOUNTS: %s", err.Error())
	}

	a.EntryDateStr = a.EntryDate.Format(common.ModelDateFormat)

	log.Printf("Inserted to SAP_ACCOUNTS: %#v", *a)

	return err
}

//
// Select some records from the resource, no use of ORP Get as it returns single record only
//
func (r *AccountRepository) ReadBulkByPartialKey(a *models.Account) (accounts []models.Account, err error) {
	log.Printf("Selecting from SAP_ACCOUNTS: %#v", *a)

	// prepeare query binding partial key value set
	var records = []models.Account{}
	var query string
	var binding map[string]interface{}
	columns := []string{
		"STATUS",
		"RELEASE_ID",
		"BSCS_ACCOUNT",
		"OFI_SAP_ACCOUNT",
		"VALID_FROM_DATE",
		"VAT_CODE_IND",
		"OFI_SAP_WBS_CODE",
		"CIT_MARKER_VAT_FLAG",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"RELEASE_DATE",
		"RELEASE_OWNER",
		"REC_VERSION",
	}
	if a.Status != "" && a.ReleaseId != "" {
		query = fmt.Sprintf(`
SELECT %s 
FROM SAP_ACCOUNTS 
WHERE STATUS = :status 
AND RELEASE_ID = :release
`, strings.Join(columns, ", "))
		binding = map[string]interface{}{
			"release": a.ReleaseId,
			"status":  a.Status,
		}
	} else {
		query = fmt.Sprintf(`
SELECT %s 
FROM SAP_ACCOUNTS Q
WHERE (STATUS IN ('W', 'C') AND RELEASE_ID = 0) 
   OR (STATUS = 'P' AND RELEASE_ID = 
      (SELECT NVL(MAX(RELEASE_ID), 0) 
         FROM SAP_ACCOUNTS SQ 
        WHERE SQ.BSCS_ACCOUNT = Q.BSCS_ACCOUNT))
`, strings.Join(columns, ", "))
	}

	// Do query by dynamicly built key
	_, err = r.Dbmap.Select(&records, query, binding)
	if err != nil {
		return nil, fmt.Errorf("Error in select from SAP_ACCOUNTS: " + err.Error())
	}

	// Take care of dates presentation
	for i, r := range records {
		if !r.ValidFromDate.IsZero() {
			records[i].ValidFromDateStr = r.ValidFromDate.Format(common.CutOffDateFormat)
		}
		if !r.EntryDate.IsZero() {
			records[i].EntryDateStr = r.EntryDate.Format(common.ModelDateFormat)
		}
		if !r.UpdateDate.IsZero() {
			records[i].UpdateDateStr = r.UpdateDate.Format(common.ModelDateFormat)
		}
		if !r.ReleaseDate.IsZero() {
			records[i].ReleaseDateStr = r.ReleaseDate.Format(common.ModelDateFormat)
		}
	}

	accounts = records

	log.Printf("Selected from SAP_ACCOUNTS records: %d %#v", len(records), records)

	return
}

//
// Update one record in resource SAP_ACCOUNTS using primary key
//
func (r *AccountRepository) UpdateByPrimaryKey(a *models.Account) (count int64, err error) {
	log.Printf("Updating SAP_ACCOUNTS: %#v", *a)

	var stmt = `
UPDATE SAP_ACCOUNTS 
SET OFI_SAP_ACCOUNT = :1, 
    VALID_FROM_DATE = :2, 
    VAT_CODE_IND = :3, 
    OFI_SAP_WBS_CODE = :4, 
    CIT_MARKER_VAT_FLAG = :5,
    UPDATE_DATE = :6, 
    UPDATE_OWNER = :7
WHERE STATUS = :8
  AND RELEASE_ID = :9
  AND BSCS_ACCOUNT = :10
`

	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(stmt,
			a.OfiSapAccount,
			a.ValidFromDate,
			a.VatCodeInd,
			a.OfiSapWbsCode,
			a.CitMarkerVatFlag,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount)
	} else {
		rs, err = r.Dbmap.Exec(stmt,
			a.OfiSapAccount,
			a.ValidFromDate,
			a.VatCodeInd,
			a.OfiSapWbsCode,
			a.CitMarkerVatFlag,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACCOUNTS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Updated SAP_ACCOUNTS records: %d", count)

	return
}

//
// Update one attribute of the record in resource SAP_ACCOUNTS using primary key
//
func (r *AccountRepository) UpdateAttributeByPrimaryKey(a *models.Account, attribute string, value interface{}) (count int64, err error) {
	log.Printf("Updating SAP_ACCOUNTS: %s <- %v %T %#v with key: %#v", attribute, value, value, value, *a)

	// make dynamic sql statement
	var stmt = `
UPDATE SAP_ACCOUNTS 
SET %s = :1, 
    UPDATE_DATE = :2, 
    UPDATE_OWNER = :3 
WHERE STATUS = :4
  AND RELEASE_ID = :5
  AND BSCS_ACCOUNT = :6
`

	tr := map[string]string{
		"status":           "STATUS",
		"releaseId":        "RELEASE_ID",
		"bscsAccount":      "BSCS_ACCOUNT",
		"ofiSapAccount":    "OFI_SAP_ACCOUNT",
		"validFromDate":    "VALID_FROM_DATE",
		"vatCodeInd":       "VAT_CODE_IND",
		"ofiSapWbsCode":    "OFI_SAP_WBS_CODE",
		"citMarkerVatFlag": "CIT_MARKER_VAT_FLAG",
		"recVersion":       "REC_VERSION",
	}

	colname, ok := tr[attribute]
	if !ok {
		return 0, fmt.Errorf("No column name for attribute: %s", attribute)
	}

	stmt = fmt.Sprintf(stmt, colname)

	// run it
	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(stmt,
			value,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount)
	} else {
		rs, err = r.Dbmap.Exec(stmt,
			value,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACCOUNTS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Updated SAP_ACCOUNTS: %s <- %v %T %#v with key: %#v, count: %d", colname, value, value, value, a, count)

	return
}

//
// Delete one record from resource SAP_ACCOUNTS using primary key
//
func (r *AccountRepository) DeleteByPrimaryKey(a *models.Account) (count int64, err error) {
	log.Printf("Deleting from SAP_ACCOUNTS: %#v", *a)

	if r.t != nil {
		count, err = r.t.Delete(a)
	} else {
		count, err = r.Dbmap.Delete(a)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Deleted from SAP_ACCOUNTS records: %d", count)

	return
}

//
// Select max
//
func (r *AccountRepository) GetMaxRelease() (release int64, err error) {
	log.Printf("Selecting MAX(RELEASE_ID) from SAP_ACCOUNTS")

	// do query
	query := "SELECT NVL(MAX(RELEASE_ID), 0) FROM SAP_ACCOUNTS"
	release, err = r.Dbmap.SelectInt(query)
	if err != nil {
		return 0, fmt.Errorf("Error in MAX(RELEASE_ID) FROM SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Selected MAX(RELEASE_ID) FROM SAP_ACCOUNTS: %d", release)

	return
}

//
// Release
//
func (r *AccountRepository) SetStatusRelease(from, into string, release, releaseNew int64) (count int64, err error) {
	log.Printf("Set STATUS, RELEASE for SAP_ACCOUNTS")

	// do update
	stmt := `
UPDATE SAP_ACCOUNTS 
SET STATUS = :1, 
    RELEASE_ID = :2 
WHERE STATUS = :3
  AND RELEASE_ID = :4
`
	var rs sql.Result
	rs, err = r.t.Exec(stmt, into, releaseNew, from, release)
	if err != nil {
		return 0, fmt.Errorf("Error in update SAP_ACCOUNTS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Updated SAP_ACCOUNTS records: %d", count)

	return
}

//
// Purge
//
func (r *AccountRepository) DeleteAll() (count int64, err error) {
	log.Printf("Purging SAP_ACCOUNTS")

	// do query
	query := "DELETE FROM SAP_ACCOUNTS"
	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(query)
	} else {
		rs, err = r.Dbmap.Exec(query)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACCOUNTS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACCOUNTS: %s", err.Error())
	}

	log.Printf("Deleted SAP_ACCOUNTS records: %d", count)

	return
}

//
// Used for validation of the lower bound of the account package
//
func (r *AccountRepository) GetMinValidDate(status string, release int64) (ts time.Time, err error) {
	query := `
SELECT NVL(MIN(VALID_FROM_DATE), SYSDATE)
  FROM SAP_ACCOUNTS 
 WHERE STATUS = :status 
   AND RELEASE_ID = :release_id`

	binding := map[string]interface{}{
		"status":     status,
		"release_id": release,
	}

	err = r.Dbmap.SelectOne(&ts, query, binding)
	if err != nil {
		err = fmt.Errorf("Error in select MIN(VALID_FROM_DATE) from SAP_ACCOUNTS: " + err.Error())
	}

	return
}

//
// Read logs of the account
//
func (r *AccountRepository) ReadLog(account string) (logs []models.AccountLog, err error) {
	var records = []models.AccountLog{}
	var query string
	var binding map[string]interface{}
	columns := []string{
		"OPCODE",
		"OPDATE",
		"STATUS",
		"BSCS_ACCOUNT",
		"OFI_SAP_ACCOUNT",
		"VALID_FROM_DATE",
		"VAT_CODE_IND",
		"OFI_SAP_WBS_CODE",
		"CIT_MARKER_VAT_FLAG",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"RELEASE_DATE",
		"RELEASE_OWNER",
		"REC_VERSION",
	}

	query = fmt.Sprintf(`SELECT %s FROM SAP_ACCOUNTS_LOG WHERE BSCS_ACCOUNT = :account`, strings.Join(columns, ","))
	binding = map[string]interface{}{
		"account": account,
	}

	_, err = r.Dbmap.Select(&records, query, binding)
	if err != nil {
		return nil, fmt.Errorf("Error in select from SAP_ACCOUNTS_LOG: " + err.Error())
	}

	// Take care of dates presentation
	for i, r := range records {
		if !r.ValidFromDate.IsZero() {
			records[i].ValidFromDateStr = r.ValidFromDate.Format(common.CutOffDateFormat)
		}
		if !r.EntryDate.IsZero() {
			records[i].EntryDateStr = r.EntryDate.Format(common.ModelDateFormat)
		}
		if !r.UpdateDate.IsZero() {
			records[i].UpdateDateStr = r.UpdateDate.Format(common.ModelDateFormat)
		}
		if !r.ReleaseDate.IsZero() {
			records[i].ReleaseDateStr = r.ReleaseDate.Format(common.ModelDateFormat)
		}
	}

	logs = records

	log.Printf("Selected from SAP_ACCOUNTS_LOG records: %d %#v", len(records), records)

	return
}
