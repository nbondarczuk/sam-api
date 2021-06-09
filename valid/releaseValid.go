package valid

import (
	"net/http"
)

func WithRelease(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
