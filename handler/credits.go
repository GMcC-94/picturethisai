package handler

import (
	"fmt"
	"net/http"
	"os"
	"picturethisai/view/credits"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return render(r, w, credits.Index())
}

func HandleStripeCallbackSuccess(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(chi.URLParam(r, "productID"))
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Change to productID
				Price:    stripe.String(""),
				Quantity: stripe.Int64(1),
			},
		},
	}

	session, err := session.New(checkoutParams)
	if err != nil {
		return err
	}

	http.Redirect(w, r, session.URL, http.StatusSeeOther)
	return nil
}
