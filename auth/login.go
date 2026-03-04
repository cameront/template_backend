package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/cameront/template_backend/logging"
	"github.com/cristalhq/jwt/v5"
)

const AuthCookieName = "ac"

type CtxKey string

const UserCtxKey CtxKey = "user"

type UserClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserAuthenticatingHandler is a handler wrapper that TODO
func UserAuthenticatingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			logging.GetLogger(ctx).Info("no cookie found")
			http.Error(w, "no cookie found", http.StatusUnauthorized)
			return
		}

		userClaims, err := validateToken(cookie.Value)
		if err != nil {
			logging.GetLogger(ctx).Info("invalid token")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		if !userClaims.IsValidAt(time.Now()) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UserCtxKey, userClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
