package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

// CRUD access metods for resource Segment
func SetDictionaryAccountBscsRoutes(router *mux.Router) *mux.Router {
	dictionaryRouter := mux.NewRouter()

	// segment access routes
	dictionaryRouter.HandleFunc("/api/dictionary/account/bscs", controllers.DictionaryAccountBscsReadAll).Methods("GET").Name("dictionary-account-bscs")

	// Handle CORS
	dictionaryRouter.HandleFunc("/api/dictionary/account/bscs", common.WithCors).Methods("OPTIONS")

	// login required before access
	router.PathPrefix("/api/dictionary/account/bscs").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithDictionaryAccountBscs),
		negroni.Wrap(dictionaryRouter),
	))

	return router
}
