package middleware

import (
	"context"
	"net/http"
	"webpolls/utils"
)

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetSession(r)
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			// Inject user info into context
			ctx := context.WithValue(r.Context(), UserIDKey, session.Values["user_id"])
			ctx = context.WithValue(ctx, UsernameKey, session.Values["username"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
