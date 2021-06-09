/*

PACKAGE: Account controller layer

It provides simplified CRUD method handlers for each API action
on the Account objects. It interracts with
data layer but it does not provide access to system
information like createion date or update date
not needed from the point of view of the frontend.

The operations On Account are:

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
func getAccountPathVars4BulkAccess(r *http.Request) (status string, release string, err error) {	
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

// Exact row lever acces by full primary key: {statu}/{release}/{account}
func getAccountPathVars4KeyAccess(r *http.Request) (status string, release string, account string, err error) {
	status, release, err = getAccountPathVars4BulkAccess(r)
	if err != nil {
		return
	}

	account, err = common.PathVariableStr(r, "account", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable account")
		return
	}

	return
}

// Exact row lever acces by partial primary key: {account}
func getAccountPathVars4LogAccess(r *http.Request) (account string, err error) {
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
func AccountCreateOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming CustomerSegment json
	var dataRequestResource resources.AccountRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Decoded payload: %#v", dataRequestResource.Data)

	// Do creation of an object
 	account := &dataRequestResource.Data
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false); 
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return		
	}
	defer repo.Close()

	if err := repo.Create(account); err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error while creating account - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return creation result with headers and appropriate status
	var dataReplyResource = resources.AccountReplyResource{Data: *account}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred", http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusCreated, j)
	}

	log.Printf("Created account, status: %d", http.StatusCreated)
}

//
// Read entitis from backend in status W or C or P if RELEASE_ID = last
//
func AccountReadActiveAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)
	
	// Perform repository parameteric read using query parameters provided
	var err error

	// Do bulk read with the key from path variables
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false);
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	// Do select on account with key pattern
	account := &models.Account{}
	accounts, err := repo.ReadBulkByPartialKey(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - " + err.Error(), http.StatusInternalServerError)
		return	
	}

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/csv":
		as := models.Accounts{Data: accounts}	
		if payload, err := as.ToCsv(); err != nil {
			common.DisplayAppError(w, common.EncoderCsvError, "Error in payload make - " + err.Error(), http.StatusInternalServerError)
			return	
		} else {
			WriteResponse(w, http.StatusOK, payload, ct)
		}
		
	case "application/xlsx":
		as := models.Accounts{Data: accounts}	
		if payload, err := as.ToExcel(); err != nil {
			common.DisplayAppError(w, common.EncoderExcelError, "Error in payload make - " + err.Error(), http.StatusInternalServerError)
			return	
		} else {
			WriteResponse(w, http.StatusNotImplemented, payload, ct)
		}
		
	default:
		// Return selection result set with headers and appropriate status
		var dataReplyResource = resources.AccountsReplyResource{
			Count: int64(len(accounts)),
			Data:  accounts,
		}
		if j, err := json.Marshal(dataReplyResource); err != nil {
			common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
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
func AccountReadSome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)
	
	// Perform repository parameteric read using query parameters provided
	var err error
	account := &models.Account{}
	if account.Status, account.ReleaseId, err = getAccountPathVars4BulkAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do bulk read with the key from path variables
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false);
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	// Obtain last release with helper
	if account.ReleaseId == "last" {
		if id, err := repo.GetMaxRelease(); err != nil {
			common.DisplayAppError(w, common.RepositoryRunError, "Error in obtaining max release - " + err.Error(), http.StatusInternalServerError)
			return
		} else {
			account.ReleaseId = fmt.Sprintf("%d", id)
		}
	}
	
	// Do select on account with key pattern
	accounts, err := repo.ReadBulkByPartialKey(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read - " + err.Error(), http.StatusInternalServerError)
		return	
	}

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.AccountsReplyResource{
		Count: int64(len(accounts)),
		Data:  accounts,
	}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}
	
	log.Printf("Read from account, status: %d", http.StatusOK)
}

//
// Makes modifications to a single record accessed
// by primary key: {statu}/{release}/{account}
//
func AccountUpdateOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Decode the incoming CustomerSegment json with values to be used in patch
	var dataRequestResource resources.AccountRequestResource
	log.Printf("Decoding payload")
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}
	account := &dataRequestResource.Data
	log.Printf("Decoded payload: %#v", account)
	
	// Decode key values of the single record to updated
	var err error
	if account.Status, account.ReleaseId, account.BscsAccount, err = getAccountPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.UpdateByPrimaryKey(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository update, no matching record found on: " + r.URL.Path, http.StatusNotFound)		
		repo.Rollback()
		return
	}

	// Return final result set with headers and appropriate status
	accounts := make([]models.Account, 1)
	accounts[0] = *account
	dataReplyResource := resources.AccountsReplyResource{Count: count, Data: accounts}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}
	
	log.Printf("Updated account, status: %d", http.StatusOK)
}

//
// Makes modifications to a record accessed
// by primary key: {statu}/{release}/{account} by specific attributes only
//
func AccountUpdateAttributes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)
	
	evals, _, err := common.GetEvaluations(r)
	if err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error in getting evaluations - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Evaluations: %#v", *evals)
	
	// Decode key values of the single record to updated
	var key = &models.Account{}
	if key.Status, key.ReleaseId, key.BscsAccount, err = getAccountPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do selective update by the composite key
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, true)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
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
			common.DisplayAppError(w, common.ControllerError, "Error in repository update, no matching record found on: " + r.URL.Path + " - attribute: " + attribute, http.StatusNotFound)
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
		}
	}

	// Make a reponse using original request body
	var dataRequestResource resources.AccountRequestResource
	if err := json.NewDecoder(r.Body).Decode(&dataRequestResource); err != nil {
		common.DisplayAppError(w, common.DecoderJsonError, "Invalid Account json request - " + err.Error(), http.StatusInternalServerError)
		return
	}

	account := &dataRequestResource.Data

	// Key values are not stored in the payload but must be taken from url path vars
	account.Status = key.Status
	account.ReleaseId = key.ReleaseId
	account.BscsAccount = key.BscsAccount
	
	// Last update timestampis used as a denominator of the whole operation
	account.UpdateDate = key.UpdateDate
	account.UpdateDateStr = key.UpdateDateStr
	account.UpdateOwner = key.UpdateOwner

	// Return final result set with headers and appropriate status
	var accounts []models.Account = make([]models.Account, 1)
	accounts[0] = *account
	dataReplyResource := resources.AccountsReplyResource{Count: count, Data: accounts}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		repo.Rollback()
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}

	repo.Commit()
	
	log.Printf("Updated account, status: %d", http.StatusOK)
}

//
// Deletes an item from the backend table,
// by primary key {release}/{status}/{account}
//
func AccountDeleteOne(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Extract parameters
	var err error
	account := &models.Account{}
	if account.Status, account.ReleaseId, account.BscsAccount, err = getAccountPathVars4KeyAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Do selective delete by the key
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	count, err := repo.DeleteByPrimaryKey(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		common.DisplayAppError(w, common.ControllerError, "Error in repository delete, no matching record found on: " + r.URL.Path, http.StatusNotFound)
		return
	}

	// Return final result set with headers and appropriate status
	var accounts []models.Account = make([]models.Account, 1)
	accounts[0] = *account	
	dataReplyResource := resources.AccountsReplyResource{Count: count, Data: accounts}
	if j, err := json.Marshal(dataReplyResource); err != nil {
		common.DisplayAppError(w, common.EncoderJsonError, "An error has occurred - " + err.Error(), http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, j)
	}
	
	log.Printf("Deleted account, status: %d", http.StatusOK)
}

//
// Purge
//
func AccountDeleteAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// Do un-selective purge
	var err error
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, err, "Error while creating account repository", http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	_, err = repo.DeleteAll()
	if err != nil {
		common.DisplayAppError(w, err, "Error in repository delete", http.StatusInternalServerError)
		return
	} else {
		WriteResponseJson(w, http.StatusOK, nil)
	}
	
	log.Printf("Purged account, status: %d", http.StatusOK)
}

//
// Read Log
//
func AccountReadLog(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	var err error
	var account string
	
	if account, err = getAccountPathVars4LogAccess(r); err != nil {
		common.DisplayAppError(w, common.ControllerError, "Error getting url variables - " + err.Error(), http.StatusInternalServerError)
		return
	}
	
	user := r.Header.Get("user")
	repo, err := repository.NewAccountRepository(user, false)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating account repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	var logs []models.AccountLog
	logs, err = repo.ReadLog(account)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository read log - " + err.Error(), http.StatusInternalServerError)
		return
	} 

	// Return selection result set with headers and appropriate status
	var dataReplyResource = resources.AccountLogsReplyResource{
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
