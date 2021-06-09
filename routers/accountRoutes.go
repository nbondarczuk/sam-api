package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

//
// Acc CRUD access metods for resource Account
//
func SetAccountRoutes(router *mux.Router) *mux.Router {
	accountRouter := mux.NewRouter()

	// account access routes
	accountRouter.HandleFunc("/api/account", controllers.AccountCreateOne).Methods("POST").Name("account")
	accountRouter.HandleFunc("/api/account", controllers.AccountReadActiveAll).Methods("GET").Name("account")
	accountRouter.HandleFunc("/api/account", controllers.AccountDeleteAll).Methods("DELETE").Name("account")
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}", controllers.AccountReadSome).Methods("GET").Name("account-status-release")
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}", controllers.AccountUpdateOne).Methods("PUT").Name("account-status-release-account")
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}", controllers.AccountUpdateAttributes).Methods("PATCH").Name("account-status-release-account")
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}", controllers.AccountDeleteOne).Methods("DELETE").Name("account-status-release-account")
	accountRouter.HandleFunc("/api/account/log/{account:[A-Za-z0-9]+}", controllers.AccountReadLog).Methods("GET").Name("account-log")

	// Handle CORS
	accountRouter.HandleFunc("/api/account/log/{account:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")	
	accountRouter.HandleFunc("/api/account/{status:[WCP]}/{release:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	accountRouter.HandleFunc("/api/account", common.WithCors).Methods("OPTIONS")
	
	// login required before access
	router.PathPrefix("/api/account").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithAccount),
		negroni.Wrap(accountRouter),
	))

	return router
}
