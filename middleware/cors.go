package middleware

import (
	"net/http"
	"maxwarden/config"
)

func EnableCors(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.DEBUG {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://"+config.GetConfig().Domain)
		}

		h.ServeHTTP(w, r)
	})
}
