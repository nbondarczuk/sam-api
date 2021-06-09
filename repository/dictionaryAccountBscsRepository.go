/*

PACKAGE: Data access layer for DictionaryAccountBscs -> GLACCOUNT table

It provides operations for accessing the data layer objects stored
in backend database. It handles the system specific attributes like
ENTRY, UPDATE dates and application used codes which are provided
by the context.

The following CRUD access methods are available:

  - ReadAll

*/

package repository

import (
	"fmt"
	"log"
	"strings"
	
	_ "gopkg.in/goracle.v2"

	"sam-api/common"
	"sam-api/models"
)

//
// Pepository being handled by request
//
type DictionaryAccountBscsRepository struct {
	Repository
}

//
// Creates new repository using existing db connection
//
func NewDictionaryAccountBscsRepository(user string) (r *DictionaryAccountBscsRepository, err error) {
	log.Printf("Creating new repository: user:" + user)

	if db, err := common.GetDbSession(); err != nil {
		return nil, err
	} else {
		dbmap := initRepository(db)
		dbmap.AddTableWithName(models.DictionaryAccountBscs{}, "GLACCOUNTS").
			SetKeys(false, "GLACODE")
		r = &DictionaryAccountBscsRepository{			
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

func (r *DictionaryAccountBscsRepository) Close() {
	r.m.Unlock()
}

//
// Select all records from the backend table, no use of ORP Get as it returns single record
//
func (r *DictionaryAccountBscsRepository) ReadAll() (entries []models.DictionaryAccountBscs, err error) {
	log.Printf("Selecting from: GLACCOUNTS")

	columns := []string{
		"GLACODE",
		"GLADESC",
		"GLATYPE",
		"GLACTIVE",
		"ENTRY_DATE",
		"ENTRY_OWNER",
		"UPDATE_DATE",
		"UPDATE_OWNER",
	}
	query := fmt.Sprintf("SELECT %s FROM GLACCOUNTS", strings.Join(columns,","))

	// do query
	records := []models.DictionaryAccountBscs{}
	_, err = r.Dbmap.Select(&records, query)
	if err != nil {
		return nil, fmt.Errorf("Error in select from GLACCOUNTS: %s", err.Error())
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

	log.Printf("Selected from GLACCOUNTS records: %d %#v", len(records), records)
	
	return
}
