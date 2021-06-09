package controllers

import (
	"log"
	"net/http"

	"sam-api/common"
)

// make a successful response with status and payload
func WriteResponseJson(w http.ResponseWriter, status int, payload []byte) {
	log.Printf("Producing payload: application/json; charset=utf-8",)		
	common.SetupCorsResponse(&w, false)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		w.Write(payload)
		log.Printf("Produced json payload len: %d, value: %s", len(payload), payload)
	} else {
		log.Printf("Produced no json payload")
	}
}

// make a successful response with status and payload
func WriteResponse(w http.ResponseWriter, status int, payload []byte, ct string) {
	log.Printf("Producing payload: %s", ct)	
	common.SetupCorsResponse(&w, false)
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(status)
	if payload != nil {
		w.Write(payload)
		log.Printf("Produced payload len: %d", len(payload))
	} else {
		log.Printf("Produced no payload")
	}
}

