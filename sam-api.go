package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/negroni"

	"sam-api/common"
	"sam-api/routers"
)

// set via loader -X flag
var (
	version, build, level string
)

// handle kill signal
func listenForShutdown(ch <-chan os.Signal) {
	<-ch
	log.Printf("Busted!")
	os.Exit(0)
}

func init() {
	common.LogInit(false)
	common.FlagsInit()
}

//
// SAP API server, may panic causing stop on init, all handlers guarded with recover
//
func main() {
	// Initialisation of the environment
	common.EnvInit(version, build, level)

	// Start all components of the server
	log.Printf("Starting API server as PID: %d in RunPath: %s", os.Getpid(), common.AppConfig.RunPath)
	log.Printf("Runing version: %s", common.GetVersion())
	common.StartUp()

	// add router ith log format change in negroni
	router := routers.InitRoutes()
	negroni.LoggerDefaultDateFormat = common.LogTimeFormat
	handler := negroni.New()
	logger := negroni.NewLogger()
	logger.SetFormat(common.LogFormat4Negroni)
	handler.Use(logger)
	handler.UseHandler(router)

	// Configure HTTP server parameters
	var server = &http.Server{
		Addr:           common.AppConfig.ServerIPAddress + ":" + common.AppConfig.ServerPort,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	// Graceful termination on SIGTERM
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		log.Printf("API server shutdown completed")
	}()

	log.Printf("API server listening on port: %s", common.AppConfig.ServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("Error API server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
