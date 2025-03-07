package app

import (
	. "maxwarden/ui"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"maxwarden/middleware"

	"net/http"
)

func VaultHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)
	session := middleware.GetSession(r)

	AppLayout("Vault", *identity, session,
		Div(
			InlineStyle(`
				$me {
					display: flex;
					flex-direction: row-reverse;
					align-items: center;
					margin-bottom: $5;
				}
			`),
			A(Href("/app/editor/add"), ButtonUI(Text("+ Add Credentials"))),
		),
		HxLoad("/app/vault-hx"),
	).Render(w)
}
