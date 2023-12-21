package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cameront/go-svelte-sqlite-template/log"
	"github.com/cristalhq/jwt/v5"
)

const authCookieKey = "ac"

type CtxKey string

const UserCtxKey CtxKey = "user"

type UserClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}

type loginRequest struct {
	Username string
	Password string
}

// LoginHandler reads a TODO
func LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.GetLogger(ctx).Info("error reading body: %v", err)
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		l := loginRequest{}
		if err = json.Unmarshal(data, &l); err != nil {
			log.GetLogger(ctx).Info("error parsing json: %v", err)
			http.Error(w, "error parsing json", http.StatusBadRequest)
			return
		}

		// Obviously you'd want to either do an oauth flow or consult the db for real user info before setting the cookie.
		if l.Username != "meuser" {
			http.Error(w, "unknown user", http.StatusInternalServerError)
			log.GetLogger(ctx).Info("unknown user")
			return
		}

		if l.Password != "pass123" {
			// ...and probably not return a different error on unknown user vs. wrong password!
			http.Error(w, "invalid username/password", http.StatusInternalServerError)
			log.GetLogger(ctx).Info("invalid password")
			return
		}

		twoDays := time.Hour * 24 * 2
		expires := time.Now().Add(twoDays)
		token, err := buildToken("123", l.Username, l.Username+"@example.com", "none", expires)
		if err != nil {
			log.GetLogger(ctx).Info("error building token: %v", err)
			http.Error(w, "error building jwt token", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{Name: authCookieKey, Value: token, Expires: expires, Path: "/"}
		http.SetCookie(w, &cookie)
		w.Write([]byte("ok"))
	})
}

// UserAuthenticatingHandler is a handler wrapper that TODO
func UserAuthenticatingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie(authCookieKey)
		if err != nil {
			log.GetLogger(ctx).Info("no cookie found")
			http.Error(w, "no cookie found", http.StatusUnauthorized)
			return
		}

		userClaims, err := validateToken(cookie.Value)
		if err != nil {
			log.GetLogger(ctx).Info("invalid token")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		if !userClaims.IsValidAt(time.Now()) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UserCtxKey, userClaims.Subject)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
