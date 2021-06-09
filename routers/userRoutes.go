package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"sam-api/common"
	"sam-api/controllers"
)

func SetUserRoutes(router *mux.Router) *mux.Router {
	userRouter := mux.NewRouter()

	// user access routes
	userRouter.HandleFunc("/api/user/login", controllers.UserLogin).Methods("POST").Name("user-login")
	userRouter.HandleFunc("/api/user/relogin", controllers.UserRelogin).Methods("POST").Name("user-relogin")
	userRouter.HandleFunc("/api/user/logoff", controllers.UserLogoff).Methods("POST").Name("user-logoff")
	userRouter.HandleFunc("/api/user/info", controllers.UserInfo).Methods("POST").Name("user-info")

	// CORS
	userRouter.HandleFunc("/api/user/login", common.WithCors).Methods("OPTIONS")
	userRouter.HandleFunc("/api/user/relogin", common.WithCors).Methods("OPTIONS")
	userRouter.HandleFunc("/api/user/logoff", common.WithCors).Methods("OPTIONS")
	userRouter.HandleFunc("/api/user/info", common.WithCors).Methods("OPTIONS")
	
	// no login requird - it is login
	router.PathPrefix("/api/user/login").Handler(negroni.New(
		negroni.HandlerFunc(common.WithLog),
		negroni.Wrap(userRouter),
	))

	// login required before
	router.PathPrefix("/api/user/relogin").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.Wrap(userRouter),
	))
	router.PathPrefix("/api/user/logoff").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.Wrap(userRouter),
	))
	router.PathPrefix("/api/user/info").Handler(negroni.New(
		negroni.HandlerFunc(common.WithAuthorize),
		negroni.HandlerFunc(common.WithLog),
		negroni.Wrap(userRouter),
	))
	
	return router
}
