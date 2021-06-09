package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	// Routes for the crud entities
	router = SetSystemRoutes(router)
	router = SetUserRoutes(router)
	router = SetAccountRoutes(router)
	router = SetReleaseRoutes(router)
	router = SetDictionaryAccountBscsRoutes(router)
	router = SetDictionaryAccountSapRoutes(router)
	router = SetOrderRoutes(router)
	router = SetDictionarySegmentRoutes(router)

	return router
}
