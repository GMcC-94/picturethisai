package handler

import (
	"net/http"
	"picturethisai/db"
	"picturethisai/pkg/kit/validate"
	"picturethisai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {

	user := getAuthenticatedUser(r)
	return render(r, w, settings.Index(user))
}

func HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	errors := settings.ProfileErrors{}
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(40)),
	}).Validate(&errors)

	if !ok {
		return render(r, w, settings.ProfileForm(params, errors))
	}

	user := getAuthenticatedUser(r)
	user.Account.Username = params.Username
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return render(r, w, settings.ProfileForm(params, settings.ProfileErrors{}))
}
