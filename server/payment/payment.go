package payment

import (
	"os"

	"github.com/stripe/stripe-go/v81"
)

func InitPayment() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = key

}
