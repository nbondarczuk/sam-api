package common

import (
	"log"
	"net/http"
	"strings"
)

// For preflight OPTIONS we must return full CORS context
// in cases of ther methods only Allow-Originn feedback expected
func SetupCorsResponse(w *http.ResponseWriter, preflight bool) {
	// all cases, even methods not OPTIONS must have Origin set
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// just OPTIONS requred
	if preflight {
		allowedMethods := []string{
			"POST",
			"GET",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		}

		allowedHeaders := []string{
			"Authorization",
			"Origin",
			"Accept",
			"X-Expires-After",
			"X-Requested-With",
			"X-Request-ID",
			"Content-Type",
			"Content-Encoding",
			"Access-Control-Allow-Headers",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		}

		(*w).Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		(*w).Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	}
}

// CORS handled as preflight with OPTIONS or in other method request contexts
func WithCors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start CORS handling on request: %#v", r)

	// preflisht request
	if (*r).Method == "OPTIONS" {
		log.Println("Received CORS preflight request")
		SetupCorsResponse(&w, true)
		w.WriteHeader(http.StatusNoContent)
	}

	log.Printf("Handled CORS preflight request, status: %d, response: %#v", http.StatusNoContent, w)
}
