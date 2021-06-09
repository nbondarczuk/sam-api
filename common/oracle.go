package common

import (
	"log"

	"database/sql"
	_ "gopkg.in/goracle.v2"
)

var (
	db *sql.DB
)

//
// Allocates db  session but no actual usage at the moment
//
func createOracleDbSession() (err error) {
	connectString := AppConfig.OracleDBUser + "/" +
		AppConfig.OracleDBPassword + "@" +
		AppConfig.OracleServiceName
	
	// open initial connection to the Oracle DB
	log.Println("Connecting to Oracle DB")
	db, err = sql.Open("goracle", connectString)
	if err != nil {
		log.Fatalf("Can not create oracle session: %s\n", err)
		return
	}

	// is connection alive?
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can not ping to Oracle: %s\n", err)
		return
	} else {
		log.Println("Ping done to Oracle DB")
	}

	return
}

//
// Once open a session then reuse it
//
func GetDbSession() (dbsession *sql.DB, err error) {
	if db == nil {
		if err = createOracleDbSession(); err != nil {
			return
		}
	} 

	dbsession = db
	
	return
}
