package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

// CRUD access metods for resource Segment
func SetDictionaryAccountSapRoutes(router *mux.Router) *mux.Router {
	dictionaryRouter := mux.NewRouter()

	// segment access routes
	dictionaryRouter.HandleFunc("/api/dictionary/account/sap", controllers.DictionaryAccountSapCreateExcel).Methods("POST").HeadersRegexp("Content-Type", "application/xlsx").Name("dictionary-account-sap")
	dictionaryRouter.HandleFunc("/api/dictionary/account/sap", controllers.DictionaryAccountSapCreate).Methods("POST").HeadersRegexp("Content-Type", "application/json").Name("dictionary-account-sap")
	dictionaryRouter.HandleFunc("/api/dictionary/account/sap", controllers.DictionaryAccountSapReadAll).Methods("GET").Name("dictionary-account-sap")
	dictionaryRouter.HandleFunc("/api/dictionary/account/sap", controllers.DictionaryAccountSapDeleteAll).Methods("DELETE").Name("dictionary-account-sap")

	// Handle CORS
	dictionaryRouter.HandleFunc("/api/dictionary/account/sap", common.WithCors).Methods("OPTIONS")

	// login required before access
	router.PathPrefix("/api/dictionary/account/sap").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithDictionaryAccountSap),
		negroni.Wrap(dictionaryRouter),
	))

	return router
}
