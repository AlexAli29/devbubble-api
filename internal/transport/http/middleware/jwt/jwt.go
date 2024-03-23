package jwt

import (
	"context"
	"devbubble-api/internal/core"
	"net/http"
)

type AuthService interface {
	CreateJWT(userId string) (string, error)

	ParseToken(tokenString string) (*core.AuthTokenClaims, error)
}

func New(authService AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Session")
			if err != nil {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, err := authService.ParseToken(cookie.Value)
			if err != nil {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "UserId", claims.UserId)

			// Call the next handler in the chain with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
