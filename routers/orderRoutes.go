package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/valid"
)

//
// Acc CRUD access metods for resource Order
//
func SetOrderRoutes(router *mux.Router) *mux.Router {
	orderRouter := mux.NewRouter()

	// order access routes
	orderRouter.HandleFunc("/api/order", controllers.OrderCreateOne).Methods("POST").Name("order")
	orderRouter.HandleFunc("/api/order", controllers.OrderReadActiveAll).Methods("GET").Name("order")	
	orderRouter.HandleFunc("/api/order", controllers.OrderDeleteAll).Methods("DELETE").Name("order")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}", controllers.OrderReadSome).Methods("GET").Name("order-status-release")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}/{segment:[A-Za-z0-9]+}", controllers.OrderUpdateOne).Methods("PUT").Name("order-status-release-account-segment")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}/{segment:[A-Za-z0-9]+}", controllers.OrderUpdateAttributes).Methods("PATCH").Name("order-status-release-account-segment")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}/{segment:[A-Za-z0-9]+}", controllers.OrderDeleteOne).Methods("DELETE").Name("order-status-release-account-segment")
	orderRouter.HandleFunc("/api/order/log/{account:[A-Za-z0-9]+}", controllers.OrderReadLog).Name("order-log")
	
	// Handle CORS
	orderRouter.HandleFunc("/api/order/log/{account:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}/{account:[A-Za-z0-9]+}/{segment:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	orderRouter.HandleFunc("/api/order/{status:[WCP]}/{release:[A-Za-z0-9]+}", common.WithCors).Methods("OPTIONS")
	orderRouter.HandleFunc("/api/order", common.WithCors).Methods("OPTIONS")

	// login required before access
	router.PathPrefix("/api/order").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.HandlerFunc(valid.WithOrder),
		negroni.Wrap(orderRouter),
	))

	return router
}
