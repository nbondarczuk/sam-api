/*

PACKAGE: Release controller layer

The operations are:

  Release

*/

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	
	"sam-api/common"
	"sam-api/repository"
)

// Exact row level acces by full primary key: {id}
func getReleasePathVars4KeyAccess(r *http.Request) (release string, err error) {
	release, err = common.PathVariableStr(r, "release", true)
	if err != nil {
		err = fmt.Errorf("Missing mandatory url path variable release")
		return
	}

	return
}

func getRelease(r *http.Request, ar *repository.AccountRepository) (release int64, err error) {
	var value string
	value, err = getReleasePathVars4KeyAccess(r)
	if err != nil {
		err = fmt.Errorf("Error while getting max: %s", err.Error())
		return
	} else if value == "last" {
		release, err = ar.GetMaxRelease()
		if err != nil {
			err = fmt.Errorf("Error in get max release of account: %s", err.Error())
			return
		}
	} else {
		release, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = fmt.Errorf("Error in parsing release: %s", err.Error())
			return
		}
	}

	return
}	

// determine role dependent transition of status
func getTransitForRole(role string) (from, into string) {
	switch role {
	case "Booker":
		from = "W"
		into = "C"

	case "Control":
		from = "C"
		into = "P"

	default:
		panic(fmt.Sprintf("Invalid role:%s", role))
	}

	return
}

//
// Change Account, Order entries status W->C or C->P depending on the role
// The value of attribute release is to be qual max(relese) of Account, Order
//
func ReleaseNew(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// determine transition depending on role, Booker: W -> C, Control: C -> P
	user := r.Header.Get("user")
	role := r.Header.Get("role")
	from, into := getTransitForRole(role) // may panic

	var err error
	var accountRepository *repository.AccountRepository
	if accountRepository, err = repository.NewAccountRepository(user, true); err != nil {
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	}
	
	// for transit W -> C no change of release
	var release, releaseNew int64
	if from == "W" {
		release = 0
		releaseNew = 0
	} else if from == "C" {
		release, err = accountRepository.GetMaxRelease()
		if err != nil {
			accountRepository.Rollback()
			common.DisplayAppError(w, err, "Error in get max release of account", http.StatusInternalServerError)
			return
		}
		releaseNew = release + 1
	}
	
	// fix account: status -> C or P, release = 0 -> max(release) + 1
	var accounts int64
	if accounts, err = accountRepository.SetStatusRelease(from, into, 0, releaseNew); err != nil {
		accountRepository.Rollback()
		common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
		return
	}

	// fix order accordingly using release of account
	var orders int64
	var orderRepository *repository.OrderRepository
	if orderRepository, err = repository.NewOrderRepository(user, true); err != nil {
		accountRepository.Rollback()
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	} else {
		if orders, err = orderRepository.SetStatusRelease(from, into, 0, releaseNew); err != nil {
			orderRepository.Rollback()
			accountRepository.Rollback()
			common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
			return
		}
	}
	
	log.Printf("Released accounts: %d, orders: %d", accounts, orders)

	// send mail if only the snmt server is known
	if common.AppConfig.AlertMailServerAddress != "" {
		log.Printf("Sending e-mail to server: %s", common.AppConfig.AlertMailServerAddress)
		err = common.MailTo(common.AppConfig.AlertMailAddress,
			fmt.Sprintf("BSCS to SAP Account/Order release done by" +
				" user: %s" +
				" role:  %s" +
				" to status: %s" +
				" of release: %d",
				user,
				role,
				into,
				releaseNew))
		if err != nil {
			orderRepository.Rollback()
			accountRepository.Rollback()
			common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
			return			
		}
	}

	WriteResponseJson(w, http.StatusOK, nil)

	accountRepository.Commit()
	orderRepository.Commit()
	
	log.Printf("Release, status: %d, response: %#v", http.StatusOK, w)
}

//
// Change Account, Order entries status W->C or C->P depending on the role
// The value of attribute release is to be qual max(relese) of Account, Order
// The release may be given in the path and not determined from last existing one.
//
func ReleaseAppend(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	// determine transition depending on role, Booker: W -> C, Control: C -> P
	user := r.Header.Get("user")
	role := r.Header.Get("role")
	from, into := getTransitForRole(role) // may panic

	var err error
	var ar *repository.AccountRepository
	if ar, err = repository.NewAccountRepository(user, true); err != nil {
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	}
	
	// for transit W -> C no change of release
	var releaseNew int64
	if from == "W" {
		releaseNew = 0
	} else if from == "C" {
		releaseNew, err = getRelease(r, ar)
		if err != nil {
			common.DisplayAppError(w, err, "Error in loading release", http.StatusInternalServerError)
			return		
		}
	}
	
	// fix account: status -> C or P, release = 0 -> release
	var accounts int64
	if accounts, err = ar.SetStatusRelease(from, into, 0, releaseNew); err != nil {
		ar.Rollback()
		common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
		return
	}

	// fix order accordingly using release of account
	var orders int64
	var or *repository.OrderRepository
	if or, err = repository.NewOrderRepository(user, true); err != nil {
		ar.Rollback()
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	} else {
		if orders, err = or.SetStatusRelease(from, into, 0, releaseNew); err != nil {
			or.Rollback()
			ar.Rollback()
			common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
			return
		}
	}
	
	log.Printf("Released accounts: %d, orders: %d", accounts, orders)

	// send mail if only the snmt server is known
	if common.AppConfig.AlertMailServerAddress != "" {
		log.Printf("Sending e-mail to server: %s", common.AppConfig.AlertMailServerAddress)
		err = common.MailTo(common.AppConfig.AlertMailAddress,
			fmt.Sprintf("BSCS to SAP Account/Order release done by" +
				" user: %s" +
				" role:  %s" +
				" to status: %s" +
				" of release: %d",
				user,
				role,
				into,
				releaseNew))
		if err != nil {
			or.Rollback()
			ar.Rollback()
			common.DisplayAppError(w, err, "Error in repository release", http.StatusInternalServerError)
			return			
		}
	}

	WriteResponseJson(w, http.StatusOK, nil)

	ar.Commit()
	or.Commit()
	
	log.Printf("Release, status: %d, response: %#v", http.StatusOK, w)
}

//
// Revoke particular release if valid date is sill in the future
// 1. check if all validDate values are in the future
// 2. check if Work is empty
// 3. move account, order entries from {status:P},{release} to {status:W},{release:0}
//
func ReleaseRevoke(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)

	var err error

	user := r.Header.Get("user")
	
	var ar *repository.AccountRepository
	if ar, err = repository.NewAccountRepository(user, true); err != nil {
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	}

	var or *repository.OrderRepository
	if or, err = repository.NewOrderRepository(user, true); err != nil {
		common.DisplayAppError(w, err, "Error while creating repository", http.StatusInternalServerError)
		return
	}
	
	release, err := getRelease(r, ar)
	if err != nil {
		common.DisplayAppError(w, err, "Error in loading release", http.StatusInternalServerError)
		return		
	}

	if _, err := ar.SetStatusRelease("P", "W", release, 0); err != nil {
		common.DisplayAppError(w, err, "Error in repository update", http.StatusInternalServerError)
		return
	}

	if _, err := or.SetStatusRelease("P", "W", release, 0); err != nil {
		ar.Rollback()
		common.DisplayAppError(w, err, "Error in repository update", http.StatusInternalServerError)
		return
	}
	
	WriteResponseJson(w, http.StatusOK, nil)

	ar.Commit()
	or.Commit()
	
	log.Printf("Release, status: %d, response: %#v", http.StatusOK, w)
}
