package middleware

import (
	"context"
	"net/http"
	"webpolls/utils"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetSession(r)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// Check if it's an HTMX request
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusOK) // HTMX expects 200 for redirect
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Inject user info into context
		ctx := context.WithValue(r.Context(), UserIDKey, session.Values["user_id"])
		ctx = context.WithValue(ctx, UsernameKey, session.Values["username"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
