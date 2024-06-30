package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"picturethisai/db"
	"picturethisai/pkg/kit/validate"
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
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}

	var errors generate.FormErrors
	ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Min(10), validate.Max(100)),
	}).Validate(&errors)
	if !ok {
		return render(r, w, generate.Form(params, errors))
	}

	fmt.Println(params)

	return nil
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
