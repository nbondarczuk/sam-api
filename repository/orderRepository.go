package repository

import (
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"
	_ "gopkg.in/goracle.v2"
	"gopkg.in/gorp.v2"

	"sam-api/common"
	"sam-api/models"
)

//
// Pepository being handled by request
//
type OrderRepository struct {
	Repository
}

//
// Create GORP context and associate structures with table name
// In fact the read method will use partial primary key as given in the route.
//
func initOrderRepository(db *sql.DB) (dbmap *gorp.DbMap) {
	dbmap = initRepository(db)

	return dbmap
}

//
// Creates new repository using existing db connection
//
func NewOrderRepository(user string, trans bool) (r *OrderRepository, err error) {
	log.Printf("Creating new repository for user: %s", user)

	if db, err := common.GetDbSession(); err != nil {
		return nil, err
	} else {
		dbmap := initRepository(db)
		dbmap.AddTableWithName(models.Order{}, "SAP_ACC_SEGM_ORDER_NUMBERS").
			SetKeys(false, "RELEASE_ID").
			SetKeys(false, "STATUS").
			SetKeys(false, "BSCS_ACCOUNT").
			SetKeys(false, "SEGMENT_CODE")
		dbmap.AddTableWithName(models.OrderLog{}, "SAP_ACC_SEGM_ORDER_NUMBERS_LOG").
			SetKeys(false, "BSCS_ACCOUNT")
		r = &OrderRepository{
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

func (r *OrderRepository) Close() {
	r.m.Unlock()
}

func (r *OrderRepository) Commit() {
	if r.t != nil {
		err := r.t.Commit()
		if err != nil {
			log.Printf("Commit error: %s", err.Error())
		}
	}
}

func (r *OrderRepository) Rollback() {
	if r.t != nil {
		err := r.t.Rollback()
		if err != nil {
			log.Printf("Rollback error: %s", err.Error())
		}
	}
}

//
// Insert new record to the resource SAP_ACC_SEGM_ORDER_NUMBERS
//
func (r *OrderRepository) Create(o *models.Order) (err error) {
	log.Printf("Inserting to SAP_ACC_SEGM_ORDER_NUMBERS: %s %#v", r.Owner, *o)

	// default value
	if o.ReleaseId == "" {
		o.ReleaseId = "0"
	}

	// default value
	if o.Status == "" {
		o.Status = "W"
	}

	if o.ValidFromDateStr != "" {
		var err error
		o.ValidFromDate, err = time.Parse(common.CutOffDateFormat, o.ValidFromDateStr)
		if err != nil {
			return err
		}
	}

	o.EntryDate = time.Now()
	o.EntryOwner = r.Owner

	if r.t != nil {
		err = r.t.Insert(o)
	} else {
		err = r.Dbmap.Insert(o)
	}
	
	if err != nil {
		return fmt.Errorf("Error in insert to SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	o.EntryDateStr = o.EntryDate.Format(common.ModelDateFormat)
	
	log.Printf("Inserted to SAP_ACC_SEGM_ORDER_NUMBERS: %#v", *o)

	return
}

//
// Select some records from the resource SAP_ACC_SEGM_ORDER_NUMBERS, no use of ORP Get  as it returns single record
//
func (r *OrderRepository) ReadBulkByPartialKey(o *models.Order) (orders []models.Order, err error) {
	log.Printf("Selecting from SAP_ACC_SEGM_ORDER_NUMBERS with: %#v", *o)

	// Prepeare query binding partial key value set
	var records = []models.Order{}
	var query string
	var binding map[string]interface{}
	columns := []string{
		"STATUS",
		"RELEASE_ID",
		"BSCS_ACCOUNT",
		"SEGMENT_CODE",
		"ORDER_NUMBER",
		"VALID_FROM_DATE",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"RELEASE_DATE",
		"RELEASE_OWNER",
		"REC_VERSION",
	}

	if o.Status != "" && o.ReleaseId != "" {
		query = fmt.Sprintf(`
SELECT %s 
  FROM SAP_ACC_SEGM_ORDER_NUMBERS 
 WHERE STATUS = :status 
   AND RELEASE_ID = :release_id`,
			strings.Join(columns, ","))
		binding = map[string]interface{}{
			"release_id": o.ReleaseId,
			"status":     o.Status,
		}
	} else {
		query = fmt.Sprintf(`
SELECT %s 
 FROM SAP_ACC_SEGM_ORDER_NUMBERS Q
WHERE (STATUS IN ('W', 'C') AND RELEASE_ID = 0) 
   OR (STATUS = 'P' AND RELEASE_ID = 
      (SELECT NVL(MAX(RELEASE_ID), 0) 
         FROM SAP_ACC_SEGM_ORDER_NUMBERS SQ 
        WHERE SQ.BSCS_ACCOUNT = Q.BSCS_ACCOUNT 
          AND SQ.SEGMENT_CODE = Q.SEGMENT_CODE))`,
			strings.Join(columns, ","))
	}

	// Do query by parameteric key
	_, err = r.Dbmap.Select(&records, query, binding)
	if err != nil {
		return nil, fmt.Errorf("Error in select from SAP_ACC_SEGM_ORDER_NUMBERS: " + err.Error())
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

	orders = records

	log.Printf("Selected SAP_ACC_SEGM_ORDER_NUMBERS records: %d %#v", len(records), records)
	
	return
}

//
// Update one record in resource SAP_ACCOUNTS using primary key
//
func (r *OrderRepository) UpdateByPrimaryKey(a *models.Order) (count int64, err error) {
	log.Printf("Updating SAP_ACC_SEGM_ORDER_NUMBERS with: %#v", *a)

	// build dynamic sql statement
	var stmt = `
UPDATE SAP_ACC_SEGM_ORDER_NUMBERS
SET VALID_FROM_DATE = :1,
	ORDER_NUMBER = :2,
	UPDATE_DATE = :3,
	UPDATE_OWNER = :4
WHERE STATUS = :5
  AND RELEASE_ID = :6
  AND BSCS_ACCOUNT = :7
  AND SEGMENT_CODE = :8
`
	// run it
	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(stmt,
			a.ValidFromDate,
			a.OrderNumber,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount,
			a.SegmentCode)
	} else {
		rs, err = r.Dbmap.Exec(stmt,
			a.ValidFromDate,
			a.OrderNumber,
			time.Now(),
			r.Owner,
			a.Status,
			a.ReleaseId,
			a.BscsAccount,
			a.SegmentCode)
	}
	
	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	log.Printf("Updated SAP_ACC_SEGM_ORDER_NUMBERS records: %d", count)
	
	return
}

//
// Update one column of the record in resource SAP_ACCOUNTS using primary key
//
func (r *OrderRepository) UpdateAttributeByPrimaryKey(o *models.Order, attribute string, value interface{}) (count int64, err error) {
	log.Printf("Updating SAP_ACC_SEGM_ORDER_NUMBERS: %s <- %v with key: %#v", attribute, value, *o)

	// build dynamic sql statement
	var stmt = `
UPDATE SAP_ACC_SEGM_ORDER_NUMBERS
SET %s = :1,
	UPDATE_DATE = :2,
	UPDATE_OWNER = :3
WHERE STATUS = :4
  AND RELEASE_ID = :5
  AND BSCS_ACCOUNT = :6
  AND SEGMENT_CODE = :7
`

	tr := map[string]string{
		"status":        "STATUS",
		"releaseId":     "RELEASE_ID",
		"bscsAccount":   "BSCS_ACCOUNT",
		"segmentCode":   "SEGMENT_CODE",
		"orderNumber":   "ORDER_NUMBER",
		"validFromDate": "VALID_FROM_DATE",
		"recVersion":    "REC_VERSION",
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
			o.Status,
			o.ReleaseId,
			o.BscsAccount,
			o.SegmentCode)
	} else {
		rs, err = r.Dbmap.Exec(stmt,
			value,
			time.Now(),
			r.Owner,
			o.Status,
			o.ReleaseId,
			o.BscsAccount,
			o.SegmentCode)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update of SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	log.Printf("Updated: SAP_ACC_SEGM_ORDER_NUMBERS records: %d", count)
	
	return
}

//
// Delete some records from resource SAP_ACC_SEGM_ORDER_NUMBERS
//
func (r *OrderRepository) DeleteByPrimaryKey(o *models.Order) (count int64, err error) {
	log.Printf("Deleting from SAP_ACC_SEGM_ORDER_NUMBERS: %#v", *o)

	// Do delete by primary key
	if r.t != nil {
		count, err = r.t.Delete(o)
	} else {
		count, err = r.Dbmap.Delete(o)
	}
	
	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	log.Printf("Deleted SAP_ACC_SEGM_ORDER_NUMBERS records: %d", count)

	return count, err
}

//
// Select max
//
func (r *OrderRepository) GetMaxRelease() (release int64, err error) {
	log.Printf("Selecting MAX(RELEASE_ID) from SAP_ACC_SEGM_ORDER_NUMBERS")

	// do query on max but table may be empty
	query := "SELECT MAX(RELEASE_ID) FROM SAP_ACC_SEGM_ORDER_NUMBERS"
	release, err = r.Dbmap.SelectInt(query)

	return
}

//
// Release
//
func (r *OrderRepository) SetStatusRelease(from, into string, release, releaseNew int64) (count int64, err error) {
	log.Printf("Set RELEASE, STATUS for SAP_ACC_SEGM_ORDER_NUMBERS")

	// do update
	stmt := `
UPDATE SAP_ACC_SEGM_ORDER_NUMBERS
SET STATUS = :1,
	RELEASE_ID = :2
WHERE STATUS = :3
AND RELEASE_ID = :4
`
	
	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(stmt, into, releaseNew, from, release)
	} else {
		rs, err = r.Dbmap.Exec(stmt, into, releaseNew, from, release)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in update SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	log.Printf("Updated SAP_ACC_SEGM_ORDER_NUMBERS records: %d", count)
	
	return
}

//
// Purge
//
func (r *OrderRepository) DeleteAll() (count int64, err error) {
	log.Printf("Deleting from SAP_ACC_SEGM_ORDER_NUMBERS")

	// do query
	query := "DELETE FROM SAP_ACC_SEGM_ORDER_NUMBERS"
	var rs sql.Result
	if r.t != nil {
		rs, err = r.t.Exec(query)
	} else {
		rs, err = r.Dbmap.Exec(query)
	}

	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in delete from SAP_ACC_SEGM_ORDER_NUMBERS: %s", err.Error())
	}
	
	log.Printf("Deleted from SAP_ACC_SEGM_ORDER_NUMBERS records: %d", count)
	
	return
}

//
// Used for validation of the lowwr bound of the order package
//
func (r *OrderRepository) GetMinValidDate(status string, release int64) (ts time.Time, err error) {
	query := `
SELECT NVL(MIN(VALID_FROM_DATE), SYSDATE)
FROM SAP_ACC_SEGM_ORDER_NUMBERS
WHERE STATUS = :status 
AND RELEASE_ID = :release`

	binding := map[string]interface{}{
		"release": release,
		"status":  status,
	}

	err = r.Dbmap.SelectOne(&ts, query, binding)
	if err != nil {
		err = fmt.Errorf("Error in select MIN(VALID_FROM_DATE) from SAP_ACCOUNTS: " + err.Error())
	}

	return
}

//
// Read logs of the account orders
//
func (r *OrderRepository) ReadLog(account string) (logs []models.OrderLog, err error) {
	var records = []models.OrderLog{}
	var query string
	var binding map[string]interface{}
	columns := []string{
		"OPCODE",
		"OPDATE",
		"STATUS",
		"RELEASE_ID",
		"BSCS_ACCOUNT",
		"SEGMENT_CODE",
		"ORDER_NUMBER",
		"VALID_FROM_DATE",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"RELEASE_DATE",
		"RELEASE_OWNER",
		"REC_VERSION",
	}

	query = fmt.Sprintf(`SELECT %s FROM SAP_ACC_SEGM_ORDER_NUMBERS_LOG WHERE BSCS_ACCOUNT = :account`, strings.Join(columns, ","))
	binding = map[string]interface{}{
		"account": account,
	}
	
	_, err = r.Dbmap.Select(&records, query, binding)
	if err != nil {
		return nil, fmt.Errorf("Error in select from SAP_ACC_SEGM_ORDER_NUMBERS_LOG: " + err.Error())
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

	log.Printf("Loaded SAP_ACC_SEGM_ORDER_NUMBERS_LOG records: %d %#v", len(records), records)

	return
}
