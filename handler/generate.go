package handler

import (
	"log/slog"
	"net/http"
	"picturethisai/db"
	"picturethisai/types"
	"picturethisai/view/generate"

	"github.com/go-chi/chi/v5"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
	}
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {

	return render(r, w, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	// Fetch image from DB
	image := types.Image{
		Status: types.ImageStatusPending,
	}
	slog.Info("checking image status:", "id", id)
	return render(r, w, generate.GalleryImage(image))
}
