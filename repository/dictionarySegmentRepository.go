/*

PACKAGE: Data access layer for DictionarySegment -> CUSTOMER_SEGMENT table

It provides operations for accessing the data layer objects stored
in backend database. It handles the system specific attributes like
ENTRY, UPDATE dates and application used codes which are provided
by the context

Values returned use the generic structures NOT containing
system specific values like ENTRY or UPDATE dates.

*/

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
type DictionarySegmentRepository struct {
	Repository
}

//
// Creates new repository using existing db connection
//
func NewDictionarySegmentRepository(user string, trans bool) (r *DictionarySegmentRepository, err error) {
	log.Printf("Creating new repository: user:" + user)

	if db, err := common.GetDbSession(); err != nil {
		return nil, err
	} else {
		dbmap := initRepository(db)
		dbmap.AddTableWithName(models.DictionarySegment{}, "CUSTOMER_SEGMENT").
			SetKeys(false, "CSTRADEREF")
		r = &DictionarySegmentRepository{
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

func (r *DictionarySegmentRepository) Close() {
	r.m.Unlock()
}

func (r *DictionarySegmentRepository) Commit() {
	if r.t != nil {
		err := r.t.Commit()
		if err != nil {
			log.Printf("Commit error: %s", err.Error())
		}
	}
}

func (r *DictionarySegmentRepository) Rollback() {
	if r.t != nil {
		err := r.t.Rollback()
		if err != nil {
			log.Printf("Rollback error: %s", err.Error())
		}
	}
}

//
// Insert new record to the backend table
//
func (r *DictionarySegmentRepository) Create(s *models.DictionarySegment) (err error) {
	log.Printf("Inserting to CUSTOMER_SEGMENT: %#v", *s)

	s.EntryDate = time.Now()
	s.EntryOwner = r.Owner

	if r.t != nil {
		err = r.t.Insert(s)
	} else {
		err = r.Dbmap.Insert(s)
	}
	
	if err != nil {
		return fmt.Errorf("Error in insert to CUSTOMER_SEGMENT: %s", err.Error())
	}

	s.EntryDateStr = s.EntryDate.Format(common.ModelDateFormat)
	
	log.Printf("Inserted to CUSTOMER_SEGMENT: %#v", *s)

	return
}

//
// Select one or all dbrecords from the backend table, no use of ORP Get as it returns single record
//
func (r *DictionarySegmentRepository) ReadAll() (segments []models.DictionarySegment, err error) {
	log.Printf("Selecting from CUSTOMER_SEGMENT")

	columns := []string{
		"CSTRADEREF",
		"SEGM_CATEGORY",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"REC_VERSION",
	}
	
	query := fmt.Sprintf("SELECT %s FROM CUSTOMER_SEGMENT", strings.Join(columns, ","))

	// do query
	records := []models.DictionarySegment{}
	_, err = r.Dbmap.Select(&records, query)
	if err != nil {
		return nil, fmt.Errorf("Error in select from CUSTOMER_SEGMENT: %s", err.Error())
	}

	// Take care of dates presentation
	for i, r := range records {
		if !r.EntryDate.IsZero() {
			records[i].EntryDateStr = r.EntryDate.Format(common.ModelDateFormat)
		}
		if !r.UpdateDate.IsZero() {
			records[i].UpdateDateStr = r.UpdateDate.Format(common.ModelDateFormat)
		}
	}

	segments = records

	log.Printf("Selected from CUSTOMER_SEGMENT records: %d %#v", len(records), records)
	
	return
}

//
// Delete all records from resource
//
func (r *DictionarySegmentRepository) DeleteAll() (count int64, err error) {
	log.Printf("Deleting from CUSTOMER_SEGMENT")

	var rs sql.Result
	rs, err = r.Dbmap.Exec("DELETE FROM CUSTOMER_SEGMENT")
	if err != nil {
		return 0, fmt.Errorf("Error delete from CUSTOMER_SEGMENT: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error delete from CUSTOMER_SEGMENT: %s", err.Error())
	}

	log.Printf("Deleted CUSTOMER_SEGMENT records: %d", count)
	
	return
}

//
// Update one record to the backend table
//
func (r *DictionarySegmentRepository) UpdateByPrimaryKey(s *models.DictionarySegment) (count int64, err error) {
	log.Printf("Updating CUSTOMER_SEGMENT: %#v", *s)

	record := *s

	record.UpdateDate = time.Now()
	record.UpdateOwner = r.Owner

	if r.t != nil {
		count, err = r.t.Update(&record)		
	} else {
		count, err = r.Dbmap.Update(&record)
	}
		
	if err != nil {
		return 0, fmt.Errorf("Error in update of CUSTOMER_SEGMENT: %s", err.Error())
	}

	log.Printf("Updated CUSTOMER_SEGMENT records: %d", count)

	return
}

//
// Delete some records from resource
//
func (r *DictionarySegmentRepository) DeleteByPrimaryKey(s *models.DictionarySegment) (count int64, err error) {
	log.Printf("Deleting from CUSTOMER_SEGMENT: %#v", *s)

	// Do delete by primary key
	if r.t != nil {
		count, err = r.t.Delete(s)
	} else {
		count, err = r.Dbmap.Delete(s)
	}
	
	if err != nil {
		return 0, fmt.Errorf("Error in delete from CUSTOMER_SEGMENT: %s", err.Error())
	}

	log.Printf("Deleted CUSTOMER_SEGMENT records: %d", count)

	return count, err
}

//
// Update one attribute of the record in resource using primary key
//
func (r *DictionarySegmentRepository) UpdateAttributeByPrimaryKey(a *models.DictionarySegment, attribute string, value interface{}) (count int64, err error) {
	log.Printf("Updating CUSTOMER_SEGMENT: %s <- %v %T %#v with key: %#v", attribute, value, value, value, *a)

	// make dynamic sql statement
	var stmt = `
UPDATE CUSTOMER_SEGMENT
SET %s = :1, 
    UPDATE_DATE = :2, 
    UPDATE_OWNER = :3 
WHERE CSTRADEREF = :4
`

	tr := map[string]string{
		"segmCategory": "SEGM_CATEGORY",
		"csTradeRef":   "CSTRADEREF",
		"revVersion":   "REC_VERSION",
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
			a.CsTradeRef)		
	} else {	
		rs, err = r.Dbmap.Exec(stmt,
			value,
			time.Now(),
			r.Owner,
			a.CsTradeRef)
	}
	
	if err != nil {
		return 0, fmt.Errorf("Error in update of CUSTOMER_SEGMENT: %s", err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error in update of CUSTOMER_SEGMENT: %s", err.Error())
	}

	log.Printf("Updated CUSTOMER_SEGMENT: %s <- %v %T %#v with key: %#v, count: %d", colname, value, value, value, a, count)	

	return
}
