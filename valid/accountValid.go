package valid

import (
	"fmt"
	"log"
	"net/http"

	"sam-api/common"
)

func WithAccount(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// skip OPTIONS, preflight and genuine requests
	if r.Method == "OPTIONS" {
		next(w, r)
		return
	}

	// all other methods

	defer func() {
        if r := recover(); r != nil {
            common.DisplayAppError(w, fmt.Errorf("Invalid request"), "Panic handler", http.StatusInternalServerError)
			return
        }
    }()

	// here goes validation of input payload
	var ok bool = true
	var info string
	role := r.Header.Get("role")
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		data, _, err := common.GetAttributesWithValues(r)
		if err != nil {
			common.DisplayAppError(w, err, "Cant get attributes of request payload", http.StatusInternalServerError)
			return
		}

		// Check for invalidation cases on particular json fields
		for k, v := range *data {
			// release id must be unsigned integer or last
			if k == "releaseId" && !validRelease(v) {
				info = "Invalid value of the field releaseId"
				ok = false
				break
			}

			switch r.Method {
			case "POST":
				// Only 0 release orders can be created
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
				
				// Control can not create accounts
				if role == "Control" {
					info = "Invalid role Control"
					ok = false;
				}
				
				// Booker can not set attribute ofiSapWbsCode
				if role == "Booker" && k == "ofiSapWbsCode" && !common.EmptyValue(v){
					info = "Invalid role Booker accessing field ofiSapWbsCode"
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
				// Booker can not set attribute ofiSapWbsCode
				if role == "Booker" {
					if common.MemberOf(k, "ofiSapWbsCode") {
						info = "Invalid role Booker accessing field ofiSapWbsCode belonging to Control"
						ok = false
					}
				}
				
				// Control can set only attribute ofiSapWbsCode or status, releaseId to move from C to W or P
				if role == "Control" {
					if !common.MemberOf(k, "ofiSapWbsCode", "status", "releaseId") {	
						info = "Invalid role Control accessing field " + k
						ok = false
					}
				}
				
				// validFromDate must be the 1st of month in future
				if k == "validFromDate" && !common.IsCutOffDate(v) {
					info = "Invalid value of validFromDate"
					ok = false
					break
				}
			}

			log.Printf("Access check: %s %s %s %#v: %v %s", role, r.Method, k, v, ok, info)
			if !ok {
				break
			}
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
