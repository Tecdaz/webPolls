package utils

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSessionStore() {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		key = "supersecretkey" // Fallback for development
	}
	Store = sessions.NewCookieStore([]byte(key))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	}
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, "webpolls-session")
	return session
}

func SaveSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) error {
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) bool {
	session := GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		return true
	}
	return false
}
