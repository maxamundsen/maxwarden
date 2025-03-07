package middleware

import (
	"context"
	"log"
	"net/http"
	"maxwarden/config"
	"maxwarden/security"
	"time"
)

type sessionKey struct{}

func LoadSession(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionMap map[string]interface{}

		sessionCookie, err := r.Cookie(config.SESSION_COOKIE_NAME)
		if err == nil {
			decryptMap, _ := security.DecryptData[map[string]interface{}]([]byte(security.DecodeBase58(sessionCookie.Value)))
			sessionMap = *decryptMap
		}

		if sessionMap == nil {
			sessionMap = make(map[string]interface{})
		}

		ctx := context.WithValue(r.Context(), sessionKey{}, sessionMap)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSession(r *http.Request) map[string]interface{} {
	session := r.Context().Value(sessionKey{}).(map[string]interface{})
	return session
}

func PutSessionCookie(w http.ResponseWriter, r *http.Request, session map[string]interface{}) {
	cookies := r.Cookies()

	// calculate total bytes used by other cookies
	var totalBytes int
	for _, cookie := range cookies {
		if cookie.Name == config.SESSION_COOKIE_NAME {
			continue
		} else {
			totalBytes += len(cookie.Value)
		}
	}

	sessionData, err := security.EncryptData(&session)
	if err != nil {
		return
	}

	length := len(sessionData) + 8

	if length+totalBytes > 4096 {
		log.Println("Attempt to generate cookie exceeding size limit for this domain")
		return
	}

	httpCookie := &http.Cookie{
		Name:     config.SESSION_COOKIE_NAME,
		Value:    security.EncodeBase58(sessionData),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, httpCookie)
}

func DeleteSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    config.SESSION_COOKIE_NAME,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})
}
