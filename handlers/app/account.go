package app

import (
	"maxwarden/middleware"
	"maxwarden/users"

	. "maxwarden/ui"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"net/http"
)

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)
	session := middleware.GetSession(r)


	type accountTableItem struct {
		Property string
		Value interface{}
	}

	cols := []string {
		"Property",
		"Value",
	}

	user, _ := users.FetchById(identity.UserID)

	entries := []accountTableItem {
		{ "UserId", user.ID },
		{ "Username", user.Username },
		{ "Name", user.Firstname + " " +  user.Lastname },
		{ "Email", user.Email },
		{ "Last Login", user.LastLogin },
	}

	func() Node {
		return AppLayout("Account", *identity, session,
			AutoTableLite(
				cols,
				entries,
				func(item accountTableItem) Node {
					return Tr(
						Td(B(Text(item.Property))),
						Td(ToText(item.Value)),
					)
				},
				AutoTableOptions{
					BorderX: true,
					Shadow: true,
				},
			),
		)
	}().Render(w)
}
