package handler

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"picturethisai/db"
	"picturethisai/pkg/sb"
	"picturethisai/types"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		session, err := store.Get(r, sessionUserKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		accessToken := session.Values[sessionAccessTokenKey]
		if accessToken == nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			slog.Error(err.Error())
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:          uuid.MustParse(resp.ID),
			Email:       resp.Email,
			LoggedIn:    true,
			AccessToken: accessToken.(string),
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)
		if !user.LoggedIn {
			path := r.URL.Path
			http.Redirect(w, r, "/login?to="+path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WithAccountSetup(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getAuthenticatedUser(r)
		account, err := db.GetAccountByUserID(user.ID)
		// User has not setup account yet
		// Redirect to /account/setup
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/account/setup", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		user.Account = account
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
