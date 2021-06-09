
package repository

import (
	"log"
	"os"
	"sync"
	
	"database/sql"
	_ "gopkg.in/goracle.v2"
	"gopkg.in/gorp.v2"
	"sam-api/common"
)

//
// Pepository being handled by request
//
type Repository struct {
	Owner string
	Db    *sql.DB
	Dbmap *gorp.DbMap
	t     *gorp.Transaction
	m     sync.RWMutex
}

//
// Create GORP context and associate structures with table name
//
func initRepository(db *sql.DB) *gorp.DbMap {
	log.Printf("Initializing repository")

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: &gorp.OracleDialect{},
	}

	if !common.TestRun {
		dbmap.TraceOn("GORP", log.New(os.Stdout, "[SAM-API] ", log.Lmicroseconds))
	}

	log.Printf("Initialized repository")

	return dbmap
}
