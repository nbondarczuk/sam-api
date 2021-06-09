package valid

import (
	"net/http"
)

func WithDictionaryAccountBscs(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
