package app

import (
	"net/http"

	"maxwarden/middleware"
	. "maxwarden/ui"
)

func GeneratorHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)
	session := middleware.GetSession(r)

	AppLayout("Generator", *identity, session,
		HxLoad("/app/generator-hx"),
	).Render(w)
}
