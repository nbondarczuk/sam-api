package valid

import (
	"net/http"
)

func IsUserValid(user string) bool {
	if len(user) > 0 {
		return true
	}
	return false
}

func IsRoleValid(role string) bool {
	if role == "Booker" || role == "Control" || role == "Admin" {
		return true
	}
	return false
}

func WithUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
