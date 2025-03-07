package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"maxwarden/auth"
	"maxwarden/config"
	"maxwarden/security"
	"maxwarden/users"
	"strings"
	"time"
)

type identityKey struct{}

func LoadIdentity(h http.HandlerFunc, requireAuth bool) http.HandlerFunc {
	loginPath := config.IDENTITY_LOGIN_PATH
	logoutPath := config.IDENTITY_LOGOUT_PATH
	defaultPath := config.IDENTITY_DEFAULT_PATH
	redirect := config.IDENTITY_AUTH_REDIRECT

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var identity *auth.Identity

		token := r.Header.Get("Authorization")

		isToken := token != ""

		// if bearer token present, use token auth, else use cookies
		if isToken {
			redirect = false
			splitToken := strings.Split(token, "Bearer ")

			if len(splitToken) >= 2 {
				token = splitToken[1]
				identity, _ = security.DecryptData[auth.Identity]([]byte(security.DecodeBase58(token)))
			}

			if identity == nil {
				blankIdentity := &auth.Identity{Authenticated: false}

				if requireAuth {
					http.Error(w, "Error: Unauthorized", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), identityKey{}, blankIdentity)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		} else {
			identityCookie, err := r.Cookie(config.IDENTITY_COOKIE_NAME)
			if err == nil {
				identity, _ = security.DecryptData[auth.Identity]([]byte(security.DecodeBase58(identityCookie.Value)))
			}

			if identity == nil {
				blankIdentity := &auth.Identity{Authenticated: false}

				if requireAuth {
					if redirect && loginPath != r.URL.Path && logoutPath != r.URL.Path {
						redirectPath := loginPath + "?redirect=" + url.QueryEscape(r.URL.String())
						http.Redirect(w, r, redirectPath, http.StatusFound)

						return
					} else if !redirect && loginPath != r.URL.Path {
						http.Error(w, "Error: Unauthorized", http.StatusUnauthorized)
						return
					}
				}

				ctx := context.WithValue(r.Context(), identityKey{}, blankIdentity)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// fetch the current user according to the database, and validate that the security stamp hasn't changed.
		// if it has, invalidate the login session.
		latestUser, _ := users.FetchById(identity.UserID)

		securityCheckFailed := latestUser.SecurityStamp != identity.SecurityStamp
		notAuthenticated := requireAuth && !identity.Authenticated
		identityExpired := identity.Expiration.Before(time.Now())

		if securityCheckFailed || notAuthenticated || identityExpired {
			if isToken {
				http.Error(w, "Error: Unauthorized", http.StatusUnauthorized)
				return
			} else {
				DeleteIdentityCookie(w, r)
				http.Redirect(w, r, loginPath, http.StatusFound)
				return
			}
		}

		if redirect && loginPath == r.URL.Path {
			http.Redirect(w, r, defaultPath, http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), identityKey{}, identity)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIdentity(r *http.Request) *auth.Identity {
	identity := r.Context().Value(identityKey{}).(*auth.Identity)
	return identity
}

func PutIdentityCookie(w http.ResponseWriter, r *http.Request, identity *auth.Identity) {
	cookies := r.Cookies()

	// calculate total bytes used by other cookies
	var totalBytes int
	for _, cookie := range cookies {
		if cookie.Name == config.IDENTITY_COOKIE_NAME {
			continue
		} else {
			totalBytes += len(cookie.Value)
		}
	}

	// A cookie serializer is a better way to handle session data. they are still
	// generated, validated, and read only by the server, but they are stored on the
	// client in a cookie.

	// For example, when a user logs into a web service, all of their auth data is
	// packed into a serialized encrypted string, which is sent via a cookie. this
	// cookie can be sent back to the page, decrypted, and de-serialized to retrieve
	// auth information in code. this is extremely fast and cheap, since you do not
	// need to store this data in a database, or even in memory.

	// Of course with this approach you must be careful not to leak the encryption
	// key, since it can be used to decrypt legitimate keys, and sign faulty ones.
	// The key should not be checked into VCS, and be regenerated if theft is
	// suspected. Resetting the key will log *everyone* out, since no sessions
	// or identities will validate.
	identityData, err := security.EncryptData(identity)
	if err != nil {
		return
	}

	length := len(identityData) + 8 // 8 additional bytes coming from somewhere ¯\_(ツ)_/¯

	if length+totalBytes > 4096 {
		log.Println("Attempt to generate cookie exceeding 4096 bytes")
		return
	}

	httpCookie := &http.Cookie{
		Name:     config.IDENTITY_COOKIE_NAME,
		Value:    security.EncodeBase58(string(identityData)),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	// if no expiry is set, cookie defaults to clear after browser closes

	http.SetCookie(w, httpCookie)
}

func DeleteIdentityCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    config.IDENTITY_COOKIE_NAME,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})
}
