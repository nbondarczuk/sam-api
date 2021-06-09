/*

PACKAGE: Data access layer for DictionaryAccountSap -> GLACCOUNT table

It provides operations for accessing the data layer objects stored
in backend database. It handles the system specific attributes like
ENTRY, UPDATE dates and application used codes which are provided
by the context.

The following CRUD access methods are available:

  - ReadAll

*/

package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "gopkg.in/goracle.v2"

	"sam-api/common"
	"sam-api/models"
)

//
// Pepository being handled by request
//
type DictionaryAccountSapRepository struct {
	Repository
}

//
// Creates new repository using existing db connection
//
func NewDictionaryAccountSapRepository(user string) (r *DictionaryAccountSapRepository, err error) {
	log.Printf("Creating new repository: user:" + user)

	if db, err := common.GetDbSession(); err != nil {
		return nil, err
	} else {	
		dbmap := initRepository(db)
		dbmap.AddTableWithName(models.DictionaryAccountSap{}, "SAP_OFI_ACCOUNTS").
			SetKeys(false, "SAP_OFI_ACCOUNT")
		r = &DictionaryAccountSapRepository{
			Repository{
				Owner: user,
				Db:    db,
				Dbmap: dbmap,
			},
		}
		r.m.Lock()
	}

	return
}

func (r *DictionaryAccountSapRepository) Close() {
	r.m.Unlock()
}

//
// Insert new record to the backend table
//
func (r *DictionaryAccountSapRepository) Create(d *models.DictionaryAccountSap) (err error) {
	log.Printf("Inserting to SAP_OFI_ACCOUNTS: %#v", *d)

	d.EntryDate = time.Now()
	d.EntryOwner = r.Owner

	err = r.Dbmap.Insert(d)
	
	if err != nil {
		return fmt.Errorf("Error in insert to SAP_OFI_ACCOUNTS: %s", err.Error())
	}

	d.EntryDateStr = d.EntryDate.Format(common.ModelDateFormat)

	log.Printf("Inserted to SAP_OFI_ACCOUNTS: %#v", *d)

	return err
}

//
// Select all records from the backend table, no use of ORP Get as it returns single record
//
func (r *DictionaryAccountSapRepository) ReadAll() (entries []models.DictionaryAccountSap, err error) {
	log.Printf("Selecting from SAP_OFI_ACCOUNTS")

	columns := []string{
		"SAP_OFI_ACCOUNT",
		"NAME",
		"STATUS",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
		"REC_VERSION",
	}	
	query := fmt.Sprintf("SELECT %s FROM SAP_OFI_ACCOUNTS", strings.Join(columns, ","))

	// do query
	records := []models.DictionaryAccountSap{}
	_, err = r.Dbmap.Select(&records, query)
	if err != nil {
		return nil, fmt.Errorf("Error in select from SAP_OFI_ACCOUNTS: " + err.Error())
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
	
	entries = records

	log.Printf("Selected from SAP_OFI_ACCOUNTS records: %d %#v", len(records), records)
	
	return
}

//
// Delete all records from resource
//
func (r *DictionaryAccountSapRepository) DeleteAll() (count int64, err error) {
	log.Printf("Deleting from: SAP_OFI_ACCOUNTS")

	var rs sql.Result
	rs, err = r.Dbmap.Exec("DELETE FROM SAP_OFI_ACCOUNTS")
	if err != nil {
		return 0, fmt.Errorf("Error delete from SAP_OFI_ACCOUNTS: " + err.Error())
	}

	count, err = rs.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Error delete from SAP_OFI_ACCOUNTS: " + err.Error())
	}

	log.Printf("Deleted SAP_OFI_ACCOUNTS records: %d", count)
	
	return
}
