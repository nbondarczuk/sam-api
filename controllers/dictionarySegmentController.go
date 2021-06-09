/*

PACKAGE: DictionarySegment controller layer

It provides simplified CRUD method handlers for each API action
on the CustomerDictionarySegment objects. It interracts with
data layer but it does not provide access to system
information like createion date or update date
not needed from the point of view of the frontend.

The operations on DictionarySegment are:

  CreateOne
  ReadAll
  DeleteAll

*/

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sam-api/common"
	"sam-api/models"
	"sam-api/repository"
	"sam-api/resources"
)

// Exact row level acces by full primary key: {id}
func getSegmentPathVars4KeyAccess(r *http.Request) (id string, err error) {
	id, err = common.PathVariableStr(r, "id", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable id")
		return
	}

	return
}

//
// Create entity in the backend, no parameters, body only used
// containing json image with values
//
func DictionarySegmentCreate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming CustomerDictionarySegment json
	var dataRequestResource resources.DictionarySegmentRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid DictionarySegment json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	segment := &dataRequestResource.Data
	log.Printf("Decoded payload: %#v", segment)
	
	// Process data sent in the request in the context of repository
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	if err := repo.Create(segment); err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error while creating segment - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Return creation result with headers and appropriate status
	log.Printf("Returning result set: %#v", *segment)
	var dataReplyResource = resources.DictionarySegmentReplyResource{Data: *segment}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusCreated, j)
	}

	log.Printf("Created segment, status: %d", http.StatusCreated)
}

//
// Read entity from backend by primary key or all of them
//
func DictionarySegmentReadAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Process data sent in the request in the context of repository
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}	
	defer repo.Close()
	
	segments, err := repo.ReadAll()
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.DictionarySegmentsReplyResource{
		Count: int64(len(segments)),
		Data:  segments,
	}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Read all segments, status: %d", http.StatusOK)
}

//
// Purge all items in the backend table
//
func DictionarySegmentDeleteAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository full delete	
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, false)
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
	dataReplyResource := resources.DictionarySegmentsReplyResource{}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Deleted all segments, status: %d", http.StatusOK)
}

//
// Update one item in the backend table
//
func DictionarySegmentUpdateOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming json
	var dataRequestResource resources.DictionarySegmentRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	segment := &dataRequestResource.Data
	log.Printf("Decoded payload: %v", segment)

	// Decode key values of the single record to updated
	var err error
	if segment.CsTradeRef, err = getSegmentPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.UpdateByPrimaryKey(segment)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository update - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository update, no matching record found on: " + r.URL.Path, http.StatusNotFound)					
		return
	}

	// Return final result set with headers and appropriate status
	segments := make([]models.DictionarySegment, 1)
	segments[0] = *segment
	dataReplyResource := resources.DictionarySegmentsReplyResource{Count: count, Data: segments}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Updated segment, status: %d", http.StatusOK)
}

//
// Update attribute of the segment
//
func DictionarySegmentUpdateAttributes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Only one can be used in request
	evals, _, err := common.GetEvaluations(r)
	if err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error in unique attribute get - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode key values of the single record to updated
	var key = &models.DictionarySegment{}
	if key.CsTradeRef, err = getSegmentPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, true)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	var count int64
	for attribute, value := range *evals {
		count, err = repo.UpdateAttributeByPrimaryKey(key, attribute, value)
		if err != nil {
			common.DisplayAppError(w, common.RepositoryRunError, "Error in repository update - " + err.Error(), http.StatusInternalServerError)
			repo.Rollback()
			return
		} else if count == 0 {
			common.DisplayAppError(w, common.ControllerError, "Error in repository update, no matching record found on: " + r.URL.Path, http.StatusNotFound)
			repo.Rollback()
			return			
		}

		// Handle key mutation
		switch attribute {
		case "csTradeRef":
			key.CsTradeRef = value.(string)
		}
	}

	// Make a reponse using original request body
	var dataRequestResource resources.DictionarySegmentRequestResource
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}

	segment := &dataRequestResource.Data

	// Key values are not stored in the payload but must be taken from url path vars
	segment.CsTradeRef = key.CsTradeRef
	
	// Last update timestampis used as a denominator of the whole operation
	segment.UpdateDate = key.UpdateDate
	segment.UpdateDateStr = key.UpdateDateStr
	segment.UpdateOwner = key.UpdateOwner
	
	// Return final result set with headers and appropriate status
	var segments []models.DictionarySegment = make([]models.DictionarySegment, 1)
	segments[0] = *segment
	dataReplyResource := resources.DictionarySegmentsReplyResource{Count: count, Data: segments}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
		repo.Rollback()
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	repo.Commit()
	
	log.Printf("Updated segment, status: %d", http.StatusOK)
}

//
// Deletes an item from the backend table by primary key {id}
//
func DictionarySegmentDeleteOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode key values of the single record to updated
	var err error
	segment := &models.DictionarySegment{}
	if segment.CsTradeRef, err = getSegmentPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do selective delete by the key
	user := r.Header.Get("user")
	repo, err := repository.NewDictionarySegmentRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.DeleteByPrimaryKey(segment)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository delete, no matching record found on: " + r.URL.Path, http.StatusNotFound)
		return
	}
	
	// Return final result set with headers and appropriate status
	var segments []models.DictionarySegment = make([]models.DictionarySegment, 1)
	segments[0] = *segment	
	dataReplyResource := resources.DictionarySegmentsReplyResource{Count: count, Data: segments}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Deleted segment, status: %d", http.StatusOK)
}
