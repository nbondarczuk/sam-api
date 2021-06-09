package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

//
// Acc CRUD access metods for resource Segment
//
func SetDictionarySegmentRoutes(router *mux.Router) *mux.Router {
	segmentRouter := mux.NewRouter()

	// segment access routes
	segmentRouter.HandleFunc("/api/dictionary/segment", controllers.DictionarySegmentCreate).Methods("POST").Name("dictionary-segment")
	segmentRouter.HandleFunc("/api/dictionary/segment", controllers.DictionarySegmentReadAll).Methods("GET").Name("dictionary-segment")
	segmentRouter.HandleFunc("/api/dictionary/segment", controllers.DictionarySegmentDeleteAll).Methods("DELETE").Name("dictionary-segment")
	segmentRouter.HandleFunc("/api/dictionary/segment/{id:[A-Za-z0-9]+}", controllers.DictionarySegmentDeleteOne).Methods("DELETE").Name("dictionary-segment-id")	
	segmentRouter.HandleFunc("/api/dictionary/segment/{id:[A-Za-z0-9]+}", controllers.DictionarySegmentUpdateOne).Methods("PUT").Name("dictionary-segment-id")
	segmentRouter.HandleFunc("/api/dictionary/segment/{id:[A-Za-z0-9]+}", controllers.DictionarySegmentUpdateAttributes).Methods("PATCH").Name("dictionary-segment-id")

	// Handle CORS
	segmentRouter.HandleFunc("/api/dictionary/segment/{id:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	segmentRouter.HandleFunc("/api/dictionary/segment", common.WithCors).Methods("OPTIONS")

	// login required before access
	router.PathPrefix("/api/dictionary/segment").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithDictionarySegment),
		negroni.Wrap(segmentRouter),
	))

	return router
}
