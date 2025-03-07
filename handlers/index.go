package handlers

import (
	. "maxwarden/ui"

	"net/http"
	"net/http/httptest"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// serve home page if route is literally '/'
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	// By default, any unmapped route will route to '/', so make sure
	// the URL is actually '/' or else 404
	if strings.HasSuffix(r.URL.Path, "/") {
		w.WriteHeader(http.StatusNotFound)

		ErrorPage(http.StatusNotFound).Render(w)
		return
	}

	rr := &httptest.ResponseRecorder{Code: http.StatusOK}

	HttpFS.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		w.WriteHeader(rr.Code)
		ErrorPage(rr.Code).Render(w)
	} else {
		HttpFS.ServeHTTP(w, r)
	}
}
