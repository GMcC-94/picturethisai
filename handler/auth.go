package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"picturethisai/db"
	"picturethisai/pkg/kit/validate"
	"picturethisai/pkg/sb"
	"picturethisai/types"
	"picturethisai/view/auth"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accessToken"
)

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.AccountSetup())
}

func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupParams{
		Username: r.FormValue("username"),
	}

	var errors auth.AccountSetupErrors
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(40)),
	}).Validate(&errors)

	if !ok {
		return render(r, w, auth.AccountSetupForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		Username: params.Username,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")
}

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	errors := auth.SignupErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password),
			validate.Message("passwords do not match"),
		),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignupForm(params, errors))
	}

	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return render(r, w, auth.SignupSuccess(user.Email))
}

func HandleResetPasswordIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.ResetPassword())
}

func HandleResetPasswordCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	params := map[string]any{
		"email":      user.Email,
		"redirectTo": "http://localhost:300/auth/reset-password",
	}

	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", sb.BaseAuthURL), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SECRET"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("supabase password recovery responded with a non 220 status: %d => %s", resp.StatusCode, string(b))
	}

	return render(r, w, auth.ResetPasswordInitiated(user.Email))

	// Using supabase is causing errors because of rate limiting
	// Potential fix is to use email provider? or custom logic to fix this.
	// if err := sb.Client.Auth.ResetPasswordForEmail(r.Context(), user.Email); err != nil {
	// return err
	// }
	// return nil
	//return render(r, w, auth.ResetPassword())
}

func HandleResetPasswordUpdate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	params := map[string]any{
		"password": r.FormValue("password"),
	}
	_, err := sb.Client.Auth.UpdateUser(r.Context(), user.AccessToken, params)
	errors := auth.ResetPasswordErrors{
		NewPassword: "Please enter a valid password",
	}
	if err != nil {
		return render(r, w, auth.ResetPasswordForm(errors))
	}

	return hxRedirect(w, r, "/")
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: os.Getenv("REDIRECT"),
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))

	}
	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")

}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}

	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = accessToken
	return session.Save(r, w)
}
