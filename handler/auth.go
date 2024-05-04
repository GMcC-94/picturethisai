package handler

import (
	"net/http"
	"picturethisai/view/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}
