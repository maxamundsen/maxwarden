package auth

import (
	. "maxwarden/ui"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"maxwarden/auth"
	"maxwarden/config"
	"maxwarden/middleware"

	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		LoginView("").Render(w)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		userid, securityStamp, authResult := auth.Authenticate(username, password)

		if !authResult {
			log.Println("Failed login attempt. Username: " + username)
			LoginView("Username or password incorrect.").Render(w)
			return
		}

		log.Println("Successful login. Username: " + username)

		// build identity info
		identity := auth.NewIdentity(userid, securityStamp, password, false)

		// serialize and send as cookie
		middleware.PutIdentityCookie(w, r, identity)

		params := r.URL.Query()
		location := params.Get("redirect")

		if len(params["redirect"]) > 0 {
			http.Redirect(w, r, location, http.StatusFound)
			return
		}

		defaultPath := config.IDENTITY_DEFAULT_PATH
		http.Redirect(w, r, defaultPath, http.StatusFound)
	}
}

func LoginView(errorMsg string) Node {
	currentYear := time.Time.Year(time.Now())

	return RootLayout("Login | MaxWarden",
		Body(
			InlineStyle(`
				$me {
					height: 100%;
					background: $color(light-grey);
				}
			`),
			Div(
				InlineStyle(`
					$me {
						display: flex;
						flex-direction: column;
						justify-content: normal;
						padding-right: $6;
						padding-left: $6;
						padding-bottom: $5;
						padding-top: $24;
					}

					@media $md {
						$me {
							padding-right: $8;
							padding-left: $8;
						}
					}
				`),
				Div(
					InlineStyle(`
						$me {
							margin-bottom: $3;
						}

						@media $sm {
							$me {
								margin-right: auto;
								margin-left: auto;
								width: 100%;
								max-width: var(--container-md);
							}
						}
					`),
					A(Href("/"),
						Img(
							InlineStyle("$me { margin-right: auto; margin-left: auto; height: $32; width: auto; }"),
							Src("/images/logo.png"),
							Alt("MaxWarden"),
						),
					),
				),
				Div(
					InlineStyle(`
						@media $sm {
							$me {
								margin-right: auto;
								margin-left: auto;
								width: 100%;
								max-width: var(--container-md);
							}
						}
					`),
					If(errorMsg != "",
						Div(
							InlineStyle(`
								$me {
									margin-top: $5;
								}

								@media $sm {
									$me {
										margin-right: auto;
										margin-left: auto;
										width: 100%;
										max-width: var(--container-md);
									}
								}
							`),
							P(InlineStyle("$me { font-size: var(--text-sm); color: $color(red-500); }"), Text(errorMsg)),
						),
					),
					Form(InlineStyle("$me { margin-top: $5; }"), Action(""), Method("POST"), AutoComplete("off"),
						H2(
							InlineStyle(`
								$me {
									margin-top: $10;
									margin-bottom: $5;
									font-weight: var(--font-weight-bold);
									font-size: var(--text-2xl);
									letter-spacing: var(--tracking-tight);
									color: $color(deep-blue);
								}
							`),
							Div(
								InlineStyle("$me { display: flex; flex-direction: row; align-items: center; gap: $4;}"),
								Icon(ICON_LOCK_KEYHOLE, 24), Text("Secure Sign In"),
							),
						),
						Div(
							Label(
								InlineStyle("$me { display: block; font-size: var(--text-xs); font-weight: var(--font-weight-normal); color: $color(neutral-700);}"),
								For("username"),
								Text("User ID"),
							),
							Div(InlineStyle("$me { margin-top: $2; }"),
								FormInput(Name("username"), Type("text"), AutoFocus(), Required()),
							),
						),
						Div(
							Label(
								InlineStyle("$me { margin-top: $5; display: block; font-size: var(--text-xs); font-weight: var(--font-weight-normal); color: $color(neutral-700);}"),
								For("password"),
								Text("Master Key"),
							),
							Div(InlineStyle("$me { margin-top: $2; }"),
								FormInput(Name("password"), Type("password"), ID("password"), Required()),
							),
						),
						Div(
							InlineStyle("$me { margin-top: $2; }"),
							Input(
								ID("show_master"),
								Type("checkbox"),
							),
							Label(
								InlineStyle("$me { margin-left: $3; display: inline; font-size: var(--text-xs); font-weight: var(--font-weight-normal); color: $color(neutral-700);}"),
								For("show_master"),
								Text("Show Master Key"),
							),
						),
						Div(
							InlineStyle(`
								$me {
									margin-top: $5;
								}
							`),
							Button(
								InlineStyle(`
									$me {
										cursor: pointer;
										width: 100%;
										padding-top: $2;
										padding-bottom: $2;
										padding-left: $5;
										padding-right: $5;
										color: $color(white);
										background-color: $color(deep-blue);
										border-radius: var(--radius-xs);
										text-align: center;
										font-size: var(--text-sm);
									}

									$me:hover {
										background-color: $color(indigo-blue);
									}
								`),
								Type("submit"),
								Text("Unlock Vault"),
							),
						),
						InlineScript(`
							let form = me();
							let btn = me("button", me());

							let check = me("#show_master");
							let passInput = me("#password");

							check.on("click", () => {
								if (passInput.type === "password") {
									passInput.type = "text";
								} else {
									passInput.type = "password";
								}
							});

							form.on("submit", () => { btn.innerHTML = "Authenticating..."; });
						`),
					),
					P(
						InlineStyle("$me { margin-top: $10; font-size: var(--text-sm); color: $color(neutral-500);}"),
						Text("© "),
						ToText(currentYear),
						Text(" Max Amundsen"),
					),
				),
			),
		),
	)
}
