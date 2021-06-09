/*

PACKAGE: DictionaryAccountSap controller layer

It provides simplified CRUD method handlers for each API action
on the DictionaryAccountSap objects. It interracts with
data layer but it does not provide access to system
information like createion date or update date
not needed from the point of view of the frontend.

The operations on DictionaryAccountSap are:

  ReadAll

*/

package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"sam-api/common"
	"sam-api/repository"
	"sam-api/resources"
)

//
// Create entity in the backend, no parameters, body only used
// containing json image with values
//
func DictionaryAccountSapCreate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming CustomerSegment json
	var dataRequestResource resources.DictionaryAccountSapRequestResource
	log.Printf("Decoding json payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Segment json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	dictionary := &dataRequestResource.Data
	log.Printf("Decoded payload: %#v", dictionary)
	
	// Do craate of the object
	user := r.Header.Get("user")
	repo, err := repository.NewDictionaryAccountSapRepository(user)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	if err := repo.Create(dictionary); err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error while creating segment - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Printf("Returning result set: %#v", *dictionary)
	}

	// Return creation result with headers and appropriate status
	var dataReplyResource = resources.DictionaryAccountSapReplyResource{Data: *dictionary}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusCreated, j)
	}

	log.Printf("Created segment, status: %d", http.StatusCreated)
}

//
// Read entity from backend all of them
//
func DictionaryAccountSapReadAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository parameteric read using query parameters provided
	user := r.Header.Get("user")
	repo, err := repository.NewDictionaryAccountSapRepository(user)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	entries, err := repo.ReadAll()
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.DictionaryAccountSapsReplyResource{
		Count: int64(len(entries)),
		Data:  entries,
	}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Read all dictionary entries, status: %d", http.StatusOK)
}

//
// Purge all items in the backend table
//
func DictionaryAccountSapDeleteAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository full delete
	user := r.Header.Get("user")
	repo, err := repository.NewDictionaryAccountSapRepository(user)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	_, err = repo.DeleteAll()
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Return final result set with headers and appropriate status
	dataReplyResource := resources.DictionaryAccountSapReplyResource{}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Deleted all segments, status: %d", http.StatusOK)
}
