package auth

import (
	"net/http"
	"maxwarden/config"
	"maxwarden/middleware"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	middleware.DeleteIdentityCookie(w, r)
	middleware.DeleteSessionCookie(w, r)

	http.Redirect(w, r, config.IDENTITY_LOGIN_PATH, http.StatusFound)
}
