package valid

import (
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"

	"sam-api/common"
)

func routeName(r *http.Request) string {
	route := mux.CurrentRoute(r)
	return route.GetName()
}

// W, C, P
func validStatus(status string) (ok bool) {
	switch status {
	case "W":
	case "C":
	case "P":
		ok = true

	default:
		ok = false
	}

	return
}

func validStatusWithPanic(r *http.Request) {
	if status, err := common.PathVariableStr(r, "status", true); err != nil {
		panic("Invalid status, missing in route")
	} else if !validStatus(status) {
		panic("Invalid status value: " + status)
	}
}

// string, last or integer >= 0
func validRelease(release interface{}) (ok bool) {
	var value string
	value, ok = release.(string)
	if !ok {
		return
	}
	
	if value == "last" {
		ok = true
	} else {
		n, err := strconv.Atoi(value)
		if err != nil {
			ok = false
		} else if n < 0 {
			ok = false
		} else {
			ok = true
		}
	}
	
	return
}
	
func validReleaseWithPanic(r *http.Request) {
	if release, err := common.PathVariableStr(r, "release", true); err != nil {
		panic("Invalid release, missing in route")
	} else if !validRelease(release) {
		panic("Invalid release value: " + release)
	}
}

// not empty
func validAccount(account string) (ok bool) {
	if len(account) > 0 {
		ok = true
	} else {
		ok = false
	}

	return
}

func validAccountWithPanic(r *http.Request) {
	if account, err := common.PathVariableStr(r, "account", true); err != nil {
		panic("Invalid account, missing in route")
	} else if !validAccount(account) {
		panic("Invalid account value: " + account)
	}
}

// not empty
func validSegment(segment string) (ok bool) {
	if len(segment) > 0 {
		ok = true
	} else {
		ok = false
	}

	return
}

func validSegmentWithPanic(r *http.Request) {
	if segment, err := common.PathVariableStr(r, "segment", true); err != nil {
		panic("Invalid segment, missing in route")
	} else if !validSegment(segment) {
		panic("Invalid segment value: " + segment)
	}
}
	
