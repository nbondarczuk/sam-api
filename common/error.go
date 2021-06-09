package common

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type (
	AppError struct {
		Error      string `json:"error"`
		Message    string `json:"message"`
		HttpStatus int    `json:"status"`
	}
	ErrorResource struct {
		Data AppError `json:"data"`
	}
)

var AuthorizationError = errors.New("Authorization error")
var ValidationError = errors.New("Payload validation error")
var DecoderJsonError = errors.New("Decoder JSON error")
var DecoderExcelError = errors.New("Decoder Excel error")
var EncoderJsonError = errors.New("Encoder JSON error")
var EncoderExcelError = errors.New("Encoder Excel error")
var EncoderCsvError = errors.New("Encoder CSV error")
var RepositoryNewError = errors.New("Repository creation error")
var RepositoryRunError = errors.New("Repository runtime error")
var ControllerError = errors.New("Controller error")

//
// Return json error feedback to the the client
//
func DisplayAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	log.Printf("Handling error: %d - %v - %s", code, handlerError, message)

	var info string
	if handlerError != nil {
		info = handlerError.Error()
	}

	var ae AppError = AppError{
		Error:      info,
		Message:    message,
		HttpStatus: code,
	}
	
	SetupCorsResponse(&w, false)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(ErrorResource{Data: ae}); err == nil {
		w.Write(j)
		log.Printf("Error handled with response: %s", string(j))
	}
}

