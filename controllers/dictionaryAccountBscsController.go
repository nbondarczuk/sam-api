/*

PACKAGE: DictionaryAccountBscs controller layer

It provides simplified CRUD method handlers for each API action
on the AccountDictionaryBscs objects. It interracts with
data layer but it does not provide access to system
information like createion date or update date
not needed from the point of view of the frontend.

The operations on DictionaryAccountBscs are:

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
// Read entity from backend all of them
//
func DictionaryAccountBscsReadAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository parameteric read using query parameters provided
	user := r.Header.Get("user")
	repo, err := repository.NewDictionaryAccountBscsRepository(user)
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
	var dataReplyResource = resources.DictionaryAccountBscssReplyResource{
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
