package handler

import (
	"log/slog"
	"net/http"
	"picturethisai/db"
	"picturethisai/types"
	"picturethisai/view/generate"
	"strconv"

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
	user := getAuthenticatedUser(r)
	prompt := "red sportscar in a garden"
	img := types.Image{
		Prompt: prompt,
		UserID: user.ID,
		Status: types.ImageStatusPending,
	}

	if err := db.CreateImage(&img); err != nil {
		return err
	}

	return render(r, w, generate.GalleryImage(img))
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	// Fetch image from DB
	image, err := db.GetImageByID(id)
	if err != nil {
		return err
	}
	slog.Info("checking image status:", "id", id)
	return render(r, w, generate.GalleryImage(image))
}
