package app

import (
	"maxwarden/entries"
	. "maxwarden/ui"

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

			var fetchErr error
			secret, fetchErr = entries.FetchSecretFromID(identity.UserID, identity.MasterKey, id)

			if fetchErr != nil {
				http.Redirect(w, r, "/app", http.StatusFound)
				return
			}
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
			ID:          r.PathValue("id"),
			Description: desc,
			URL:         url,
			Notes:       notes,
			Password:    password,
			Username:    username,
		}

		if editorType == EDITOR_TYPE_ADD {
			entries.Add(identity.UserID, identity.MasterKey, secret)
		} else {
			entries.Update(identity.UserID, identity.MasterKey, secret)
		}

		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}

	AppLayout(title, *identity, session,
		Modal(
			"password_generator",
			nil,
			HxLoad("/app/generator-hx"),
			[]Node {
				ButtonUIOutline(ModalCloser(), Text("Close")),
			},
		),

		If(editorType == EDITOR_TYPE_EDIT,
			Group{
				Modal(
					"warning_popup",
					Text("Warning!"),
					Text("Are you sure you want to delete this entry? This action cannot be undone."),
					[]Node{
						A(Href("/app/delete/"+secret.ID), ButtonUIDanger(Text("Delete"))),
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
			FormInput(Type("text"), Name("description"), Value(secret.Description), AutoFocus(), Required()),
			Br(),

			FormLabel(Text("Username")),
			FormInput(Type("text"), Name("un"), Value(secret.Username)),
			Br(),

			Div(
				FlexLeftRight(
					FormLabel(Text("Password")),
					ModalActuator("password_generator", ButtonUI(Type("button"), Text("Generate a secure password"))),
				),
				Br(),

				FormInput(Class("password"), Type("password"), Name("pas"), Value(secret.Password)),
				Br(),

				ButtonUIOutline(Class("passtoggle"),
					Type("button"),
					Flex(Icon(ICON_EYE, 24), Text(" Toggle password visibility")),
				),

				InlineScript(`
					let toggle = me(".passtoggle", me());
					let passInput = me(".password", me());

					toggle.on("click", () => {
						if (passInput.type === "password") {
							passInput.type = "text";
						} else if (passInput.type === "text") {
							passInput.type = "password";
						}
					});
				`),
			),

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

			If(editorType == EDITOR_TYPE_EDIT,
				Div(
					InlineStyle("$me { margin-top: $3; color: $color(neutral-400); font-size: var(--text-sm);}"),
					Text("Last modified: "), FormatDateTime(secret.Modified),
					Br(),
					Text("Created: "), FormatDateTime(secret.Created),
				),
			),
		),
	).Render(w)
}
