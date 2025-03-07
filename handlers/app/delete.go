package app

import (
	"maxwarden/entries"
	"maxwarden/middleware"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)

	id := r.PathValue("id")

	entries.DeleteSecret(identity.UserID, identity.MasterKey, id)

	http.Redirect(w, r, "/app", http.StatusFound)
}