package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

// CRUD access metods for resource Release
func SetReleaseRoutes(router *mux.Router) *mux.Router {
	releaseRouter := mux.NewRouter()

	// segment access routes
	releaseRouter.HandleFunc("/api/release/new", controllers.ReleaseNew).Methods("POST").Name("release-new")
	releaseRouter.HandleFunc("/api/release/{release:[A-Za-z0-9]+}", controllers.ReleaseAppend).Methods("POST").Name("release-id")
	releaseRouter.HandleFunc("/api/release/{release:[A-Za-z0-9]+}", controllers.ReleaseRevoke).Methods("DELETE").Name("release-id")	

	// handle CORS
	releaseRouter.HandleFunc("/api/release/{release:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	releaseRouter.HandleFunc("/api/release/new", common.WithCors).Methods("OPTIONS")

	// login required before access
	router.PathPrefix("/api/release").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithRelease),
		negroni.Wrap(releaseRouter),
	))

	return router
}
