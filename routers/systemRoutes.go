package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
)

//
// CRUD access metods for resource System
//
func SetSystemRoutes(router *mux.Router) *mux.Router {
	systemRouter := mux.NewRouter()

	// segment access routesy
	systemRouter.HandleFunc("/api/system/version", controllers.SystemVersionRead).Methods("GET").Name("system-version")
	systemRouter.HandleFunc("/api/system/stat", controllers.SystemStatRead).Methods("GET").Name("system-stat")
	systemRouter.HandleFunc("/api/system/health", controllers.SystemHealthRead).Methods("GET").Name("system-health")

	// cors on all the paths
	systemRouter.HandleFunc("/api/system/version", common.WithCors).Methods("OPTIONS")
	systemRouter.HandleFunc("/api/system/stat", common.WithCors).Methods("OPTIONS")
	systemRouter.HandleFunc("/api/system/health", common.WithCors).Methods("OPTIONS")

	router.PathPrefix("/api/system").Handler(negroni.New(
		negroni.HandlerFunc(common.WithLog),
		negroni.Wrap(systemRouter),
	))

	return router
}
