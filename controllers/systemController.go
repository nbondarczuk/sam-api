package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"sam-api/common"
	"sam-api/models"
	"sam-api/resources"
)

//
// version
//
func SystemVersionRead(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Write payload to the response
	version := &models.Version{Version: common.GetVersion()}
	dataReplyResource := resources.VersionResource{Status: "Ok", Data: *version}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, err, "Error json encoding version info", http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Done system version, status: %d, response: %#v", http.StatusOK, w)
}

//
// memory stat
//
func SystemStatRead(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Write payload to the response
	stat := common.NewStat()
	dataReplyResource := resources.StatResource{Status: "Ok", Data: *stat}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, err, "Error json encoding version info", http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Done system stat, status: %d, response: %#v", http.StatusOK, w)
}

//
// needed by kubernetes?
//
func SystemHealthRead(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	WriteResponseJson(w, http.StatusOK, nil)

	log.Printf("Done system stat, status: %d, response: %#v", http.StatusOK, w)
}
