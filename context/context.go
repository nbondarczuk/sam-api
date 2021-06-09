package context

import (
	"database/sql"
	_ "gopkg.in/goracle.v2"
	"gopkg.in/gorp.v2"

	"sam-api/common"
)

//
// Struct used for maintaining HTTP Request Context
//
type Context struct {
	User  string
	Role  string
	Db    *sql.DB
	Dbmap *gorp.DbMap
}

//
// Close the session returning it to the pool
//
func (c *Context) Close() {
	c.Db.Close()
}

//
// Create a new Context object for each HTTP request using
// provided credentials
//
func NewContext(user, role string) (context *Context, err error) {
	if db, err := common.GetDbSession(); err == nil { 
		context = &Context{
			User:  user,
			Role:  role,
			Db:    db,
			Dbmap: &gorp.DbMap{Db: db, Dialect: gorp.OracleDialect{}},
		}
	}
	
	return
}
