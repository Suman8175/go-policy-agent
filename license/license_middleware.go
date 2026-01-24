package license

import (
	"fmt"
	"log"
	"net/http"
)

// Middleware wraps an HTTP handler and enforces license restrictions
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		query := r.URL.RawQuery

		fullPath := path
		if query != "" {
			fullPath = fmt.Sprintf("%s?%s", path, query)
		}

		log.Printf("Checking license status for HTTP Method=[%s] Path=[%s]...", method, fullPath)

		status := GetStatus()

		// Enforcement logic
		if status == OK || (status == SOFT_LOCK && method == http.MethodGet) {
			next.ServeHTTP(w, r)
			return
		}

		// License violation
		errMsg := fmt.Sprintf("License restriction in effect: %s", status)
		log.Printf("License violation: %s", errMsg)
		http.Error(w, errMsg, http.StatusForbidden)
	})
}
