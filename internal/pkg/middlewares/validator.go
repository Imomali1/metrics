package middlewares

import (
	"net/http"
	"strings"
)

func ValidateUpdateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		path := r.URL.RequestURI()
		path = strings.Trim(path, "/")
		params := strings.Split(path, "/")
		if len(params) != 4 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
