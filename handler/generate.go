package handler

import (
	"net/http"
	"picturethisai/view/generate"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, generate.Index())
}
