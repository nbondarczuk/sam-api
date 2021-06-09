package valid

import (
	"net/http"
)

func WithDictionaryAccountSap(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
