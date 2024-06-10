package handler

import (
	"context"
	"fmt"
	"net/http"
	"picturethisai/types"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			fmt.Println("from the user middleware")
			next.ServeHTTP(w, r)

		}
		user := types.AuthenticatedUser{}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(fn)
}
