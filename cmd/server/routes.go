package main

import (
	// "maxwarden/middleware"
	"net/http"
	"maxwarden/handlers"
	"maxwarden/handlers/app"
	"maxwarden/handlers/auth"
	"maxwarden/middleware"
)

func mapRoutes(mux *http.ServeMux) {
	id := middleware.LoadIdentity
	sess := middleware.LoadSession
	// cors := middleware.EnableCors

	mux.HandleFunc("/app/account", id(sess(app.AccountHandler), true))
	mux.HandleFunc("/app", id(sess(app.VaultHandler), true))
	mux.HandleFunc("/app/generator", id(sess(app.GeneratorHandler), true))
	mux.HandleFunc("/app/generator-hx", id(sess(app.GeneratorHxHandler), true))
	mux.HandleFunc("/app/vault-hx", id(sess(app.VaultHxHandler), true))
	mux.HandleFunc("/app/delete/{id}", id(sess(app.DeleteHandler), true))
	mux.HandleFunc("/app/editor/add", id(sess(app.EditorHandler), true))
	mux.HandleFunc("/app/editor/edit/{id}", id(sess(app.EditorHandler), true))
	mux.HandleFunc("/auth/login", id(sess(auth.LoginHandler), true))
	mux.HandleFunc("/auth/logout", id(sess(auth.LogoutHandler), true))
	mux.HandleFunc("/", handlers.IndexHandler)
}
