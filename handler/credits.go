package handler

import (
	"net/http"
	"picturethisai/view/credits"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return render(r, w, credits.Index())
}
