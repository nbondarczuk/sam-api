package valid

import (
	"fmt"
	"log"
	"net/http"
	
	"sam-api/common"
)

func WithOrder(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// skip OPTIONS, preflight and genuine requests
	if r.Method == "OPTIONS" {
		next(w, r)
		return
	}

	// all other methods
	
	defer func() {
        if r := recover(); r != nil {
			common.DisplayAppError(w, fmt.Errorf("Invalid request"), "Validation handler", http.StatusInternalServerError)
			return
        }
    }()
	
	// Check for invalidation cases on particular json fields
	var ok bool = true
	var info string
	role := r.Header.Get("role")
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		data, _, err := common.GetAttributesWithValues(r)
		if err != nil {
			common.DisplayAppError(w, err, "Cant get attributes of request payload", http.StatusInternalServerError)
			return
		}

		for k, v := range *data {
			// release id must be unsigned integer or last			
			if k == "releaseId" && !validRelease(v) {
				info = "Invalid value of the field releaseId"
				ok = false
				break
			}
			
			switch r.Method {
			case "POST":
				// Only 0 release orders can be create
				if k == "releaseId" && v != "0" {
					info = "Invalid value of release, only 0 allowed"
					ok = false
					break
				}

				// Only orders in W status can be created
				if k == "status" && v != "W" {
					info = "Invalid value of status, only W allowed"
					ok = false
					break
				}
				
				// Booker can not set attribute orderNmber
				if role == "Booker" && k == "orderNmber" && !common.EmptyValue(v){
					info = "Invalid role Booker accessing " + k
					ok = false
				}
				
				// validFromDate must be the 1st of month in future
				if k == "validFromDate" && !common.IsCutOffDate(v) {
					info = "Invalid value of validFromDate"
					ok = false
					break
				}				
				
			case "PUT": fallthrough
			case "PATCH":
				// Booker can not set attribute orderNumber
				if role == "Booker" && k == "orderNumber" {
					info = "Invalid role Booker accessing " + k
					ok = false
				}
				
				// Control can set only one attribute orderNumber
				if role == "Control" && k != "orderNumber" {
					info = "Invalid role Booker accessing " + k
					ok = false
				}

				// validFromDate must be the 1st of month in future
				if k == "validFromDate" && !common.IsCutOffDate(v) {
					info = "Invalid value of validFromDate"
					ok = false
					break
				}				
			}
			
			log.Printf("Access check: %s %s %s %#v: %v", role, r.Method, k, v, ok)
			if !ok {
				break
			}
		}
	} else if r.Method == "DELETE" {
		// Booker can not delete
		if role == "Booker" {
			info = "Invalid role Booker"
			log.Printf("Access check: %s %s: %v %s", role, r.Method, ok, info)
			ok = false
		}
	}

	if !ok {
		log.Printf("Validation error: " + info)
		common.DisplayAppError(w, common.ValidationError, info, http.StatusForbidden)
		return
	}

	log.Printf("Validation status: %v", ok)
	
	next(w, r)
}
