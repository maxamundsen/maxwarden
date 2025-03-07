package app

import (
	"maxwarden/entries"
	"maxwarden/security"
	. "maxwarden/ui"
	"maxwarden/users"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"maxwarden/middleware"
	"net/http"
)

const (
	EDITOR_TYPE_EDIT = iota
	EDITOR_TYPE_ADD  = iota
)

func EditorHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)
	session := middleware.GetSession(r)

	var editorType int
	var title string
	var btnLabel string

	if r.URL.Path == "/app/editor/add" {
		editorType = EDITOR_TYPE_ADD
		title = "Add Credentials"
		btnLabel = "Add"
	} else {
		editorType = EDITOR_TYPE_EDIT
		title = "Edit Credentials"
		btnLabel = "Save"
	}

	var secret entries.Secret

	if r.Method == http.MethodGet {
		if editorType == EDITOR_TYPE_EDIT {
			id := r.PathValue("id")

			secret, _ = entries.FetchSecretFromID(identity.UserID, identity.MasterKey, id)
		}
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		desc := r.FormValue("description")
		notes := r.FormValue("notes")
		username := r.FormValue("un")
		password := r.FormValue("pas")
		url := r.FormValue("url")

		secret = entries.Secret{
			Description: desc,
			URL:         url,
			Notes:       notes,
			Password:    password,
			Username:    username,
		}

		user, _ := users.FetchById(identity.UserID)

		// Get current secret store
		secrets, _ := security.DecryptDataWithKey[[]entries.Secret](user.Data, identity.MasterKey)

		if secrets == nil {
			http.Redirect(w, r, "/app", http.StatusFound)
			return
		}

		if editorType == EDITOR_TYPE_ADD {
			secret.ID = security.RandBase58String(32)
			*secrets = append(*secrets, secret)
		} else {
			secret.ID = r.PathValue("id")

			// linear search and replace
			for i, v := range *secrets {
				if v.ID == secret.ID {
					(*secrets)[i] = secret
				}
			}
		}

		// Serialize and encrypt modified store using master key
		enc, _ := security.EncryptDataWithKey(secrets, identity.MasterKey)

		user.Data = enc

		users.Update(user)

		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}

	AppLayout(title, *identity, session,
		If(editorType == EDITOR_TYPE_EDIT,
			Group{
				Modal(
					"warning_popup",
					Text("Warning!"),
					Text("Are you sure you want to delete this entry? This action cannot be undone."),
					[]Node{
						A(Href("/app/delete/" + secret.ID), ButtonUIDanger(Text("Delete"))),
						ButtonUIOutline(ModalCloser(), Text("Close")),
					},
				),
				Div(
					InlineStyle("$me { display: flex; flex-direction: row-reverse; align-items: center; }"),
					ModalActuator("warning_popup", ButtonUIDanger(Text("Delete"))),
				),
			},
		),
		Form(
			AutoComplete("off"),
			Method("POST"),
			FormLabel(Text("Description")),
			FormInput(Type("text"), Name("description"), Value(secret.Description)),
			Br(),

			FormLabel(Text("Username")),
			FormInput(Type("text"), Name("un"), Value(secret.Username)),
			Br(),

			FormLabel(Text("Password")),
			FormInput(Type("password"), Name("pas"), Value(secret.Password)),
			Br(),

			FormLabel(Text("URL")),
			FormInput(Type("text"), Name("url"), Value(secret.URL)),
			Br(),

			FormLabel(Text("Additional Notes")),
			FormTextarea(InlineStyle("$me { height: $32; font-family: var(--font-mono); }"), Name("notes"), Text(secret.Notes)),
			Br(),

			Div(
				InlineStyle("$me { display: flex; flex-direction: row; align-items: center; gap: $4; }"),
				ButtonUISuccess(Text(btnLabel), Type("submit")),
				A(Href("/app"), ButtonUIOutline(Text("Close"), Type("button"))),
			),
		),
	).Render(w)
}
