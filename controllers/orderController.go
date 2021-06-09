/*

PACKAGE: Order controller layer

It provides simplified CRUD method handlers for each API action
on the Account objects. It interracts with
data layer but it does not provide access to system
information like createion date or update date
not needed from the point of view of the frontend.

The operations on Order are:

  CreateOne
  ReadSome
  UpdateOne
  DeleteOne

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

// Bulk access (many rows) by partual key: {status}/{release}
func getOrderPathVars4BulkAccess(r *http.Request) (status, release string, err error) {
	status, err = common.PathVariableStr(r, "status", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable status")
		return
	}

	release, err = common.PathVariableStr(r, "release", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable release")
		return
	}

	return
}

// Exact row lever acces by full primary key: {statu}/{release}/{account}/{segment}
func getOrderPathVars4KeyAccess(r *http.Request) (status, release, account, segment string, err error) {
	status, release, err = getAccountPathVars4BulkAccess(r)
	if err != nil {
		return
	}

	account, err = common.PathVariableStr(r, "account", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable account")
		return
	}

	segment, err = common.PathVariableStr(r, "segment", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable segment")
		return
	}

	return
}

// Exact row lever acces by partial primary key: {account}
func getOrderPathVars4LogAccess(r *http.Request) (account string, err error) {
	account, err = common.PathVariableStr(r, "account", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable account")
		return
	}

	return
}

//
// Create entity in the backend, no parameters,
// body only used containing json image with values
//
func OrderCreateOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Process data sent in the request in the context of repository

	// Decode the incoming Order json
	var dataRequestResource resources.OrderRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Order json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	order := &dataRequestResource.Data
	log.Printf("Decoded payload: %v", order)
	

	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()

	// Do create object
	if err := repo.Create(order); err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error while creating order - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Return creation result with headers and appropriate status
	var dataReplyResource = resources.OrderReplyResource{Data: *order}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
	} else {
		WriteResponseJson(w, http.StatusCreated, j)
	}

	log.Printf("Created order, status: %d", http.StatusCreated)
}

//
// Read entity from backend in status W or C or P if RELEASE_ID = last
//
func OrderReadActiveAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository parameteric read using query parameters provided
	var err error

	// do bulk read
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	// Do select on account with key pattern
	order := &models.Order{}	
	orders, err := repo.ReadBulkByPartialKey(order)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - " + err.Error(), http.StatusInternalServerError)
		return
	}

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/csv":
		os := models.Orders{Data: orders}
		if payload, err := os.ToCsv(); err != nil {
			common.DisplayAppError(w, common.EncoderCsvError, "Error in payload make - " + err.Error(), http.StatusInternalServerError)
			return	
		} else {
			WriteResponse(w, http.StatusOK, payload, ct)
		}
		
	case "application/xlsx":
		os := models.Orders{Data: orders}
		if payload, err := os.ToExcel(); err != nil {
			common.DisplayAppError(w, common.EncoderExcelError, "Error in payload make - " + err.Error(), http.StatusInternalServerError)
			return	
		} else {			
			WriteResponse(w, http.StatusNotImplemented, payload, ct)
		}
		
	default:
		// Return selection result set with headers and appropriate status
		var dataReplyResource = resources.OrdersReplyResource{
			Count: int64(len(orders)),
			Data:  orders,
		}
		if j, err := json.Marshal(dataReplyResource); err != nil {
			common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
			return
		} else {
			WriteResponseJson(w, http.StatusOK, j)
		}
	}
	
	log.Printf("Read from account, status: %d", http.StatusOK)
}

//
// Read entity from backend by the part
// of primary key: {status}/{release} - mandatory parameters
//
func OrderReadSome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Perform repository parameteric read using query parameters provided
	var err error
	order := &models.Order{}
	if order.Status, order.ReleaseId, err = getOrderPathVars4BulkAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// do bulk read
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating order repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	// Obtain last release with helper
	if order.ReleaseId == "last" {
		if id, err := repo.GetMaxRelease(); err != nil {
			common.DisplayAppError(w, common.RepositoryRunError, "Error in obtaining max release - " + err.Error(), http.StatusInternalServerError)
			return
		} else {
			order.ReleaseId = fmt.Sprintf("%d", id)
		}
	}

	// Do select on account with key pattern
	orders, err := repo.ReadBulkByPartialKey(order)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - "+ err.Error(), http.StatusInternalServerError)
		return
	}

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.OrdersReplyResource{
		Count: int64(len(orders)),
		Data:  orders,
	}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Read from account, status: %d", http.StatusOK)
}

//
// Makes modifications to a single record accessed
// by primary key: {status}/{release}/{account}/{segment}
//
func OrderUpdateOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming Order json
	var dataRequestResource resources.OrderRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	order := &dataRequestResource.Data
	log.Printf("Decoded payload: %v", order)

	// Decode key values of the single record to updated
	var err error
	if order.Status, order.ReleaseId, order.BscsAccount, order.SegmentCode, err = getOrderPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.UpdateByPrimaryKey(order)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository update, no matching record found on: " + r.URL.Path, http.StatusNotFound)
		repo.Rollback()
		return
	}
	
	// Return final result set with headers and appropriate status
	orders := make([]models.Order, 1)
	orders[0] = *order
	dataReplyResource := resources.OrdersReplyResource{Count: count, Data: orders}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Updated order, status: %d", http.StatusOK)
}

//
// Makes modifications to a single record accessed
// by primary key: {status}/{release}/{account}/{segment} by specific attributes
//
func OrderUpdateAttributes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	evals, _, err := common.GetEvaluations(r)
	if err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error in getting evaluations - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Evaluations: %#v", *evals)

	// Decode key values of the single record to updated
	var key = &models.Order{}
	if key.Status, key.ReleaseId, key.BscsAccount, key.SegmentCode, err = getOrderPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, true)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository- " + err.Error(), http.StatusInternalServerError)
		return
	}

	var count int64
	for attribute, value := range *evals {
		count, err = repo.UpdateAttributeByPrimaryKey(key, attribute, value)
		if err != nil {
			common.DisplayAppError(w, common.RepositoryRunError, "Error in repository update - " + err.Error(), http.StatusInternalServerError)
			repo.Rollback()
			return
		} else if count == 0 {
			common.DisplayAppError(w, common.ControllerError, "Error in repository delete, no matching record found on: " + r.URL.Path + " - attribute: " + attribute, http.StatusNotFound)					
			repo.Rollback()
			return
		}
		
		// Handle key mutation
		switch attribute {
		case "status":
			key.Status = value.(string)
		case "releaseId":
			key.ReleaseId = value.(string)
		case "bscsAccount":
			key.BscsAccount = value.(string)
		case "segmentCode":
			key.SegmentCode = value.(string)
		}		
	}

	// Make a reponse using original request body
	var dataRequestResource resources.OrderRequestResource
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}

	order := &dataRequestResource.Data

	// Key values are not stored in the payload but must be taken from url path vars
	order.Status = key.Status
	order.ReleaseId = key.ReleaseId
	order.BscsAccount = key.BscsAccount
	order.SegmentCode = key.SegmentCode
		
	// Last update timestampis used as a denominator of the whole operation
	order.UpdateDate = key.UpdateDate
	order.UpdateDateStr = key.UpdateDateStr
	order.UpdateOwner = key.UpdateOwner
		
	// Return final result set with headers and appropriate status
	var orders []models.Order = make([]models.Order, 1)
	orders[0] = *order
	dataReplyResource := resources.OrdersReplyResource{Count: count, Data: orders}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
		repo.Rollback()
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	repo.Commit()
	
	log.Printf("Updated order, status: %d", http.StatusOK)
}

//
// Deletes an item from the backend table,
// by primary key {release}/{status}/{account}/{segment}
//
func OrderDeleteOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode key values of the single record to updated
	var err error
	order := &models.Order{}
	if order.Status, order.ReleaseId, order.BscsAccount, order.SegmentCode, err = getOrderPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do selective delete by the key
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.DeleteByPrimaryKey(order)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository delete, no matching record found on: " + r.URL.Path, http.StatusNotFound)		
		return
	}
	
	// Return final result set with headers and appropriate status
	var orders []models.Order = make([]models.Order, 1)
	orders[0] = *order	
	dataReplyResource := resources.OrdersReplyResource{Count: count, Data: orders}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An unexpected error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Deleted order, status: %d", http.StatusOK)
}

// Purge all items
func OrderDeleteAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Do un-selective purge
	var err error
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	_, err = repo.DeleteAll()
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, nil)
	}
	
	log.Printf("Purged account, status: %d", http.StatusOK)
}

//
// Read Log
//
func OrderReadLog(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	var err error
	var account string
	
	if account, err = getOrderPathVars4LogAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	user := r.Header.Get("user")
	repo, err := repository.NewOrderRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	var logs []models.OrderLog
	logs, err = repo.ReadLog(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error in repository read log - " + err.Error(), http.StatusInternalServerError)
		return
	} 

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.OrderLogsReplyResource{
		Count: int64(len(logs)),
		Data:  logs,
	}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	log.Printf("Read log for account, status: %d", http.StatusOK)
}
