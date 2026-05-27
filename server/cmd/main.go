package main

import (
	"encoding/json"
	"fmt"
	"io"
	"license-server/db"
	"license-server/internals/handler"
	"license-server/internals/middleware"
	"license-server/internals/repository"
	"license-server/internals/service"
	"license-server/payment"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {

	db := db.Connect()
	payment.InitPayment()
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewHandler(userService)

	licenseRepo := repository.NewLicenseRepository(db)
	licenseService := service.NewLicenseService(licenseRepo)
	licenseHandler := handler.NewLicenseHandler(licenseService, userService)

	paymentService := service.NewPaymentService()
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signin", userHandler.SigninHandler)
		r.Post("/signup", userHandler.SignupHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Route("/license", func(r chi.Router) {
			r.Post("/generate", licenseHandler.GenerateLicense)
			r.Post("/activate", licenseHandler.ActivateLicense)
			r.Post("/decode", licenseHandler.DecodeLicense)
		})

		r.Route("/payment", func(r chi.Router) {
			r.Post("/stripe", paymentHandler.CreateCheckout)
			r.Post("/stripeWebhook", func(w http.ResponseWriter, r *http.Request) {

				signature := r.Header.Get("stripe-signature")
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "error reading body", http.StatusBadRequest)
					return
				}

				// Verify the webhook signature
				webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
				event, err := webhook.ConstructEvent(body, signature, webhookSecret)
				if err != nil {
					http.Error(w, "invalid signature", http.StatusBadRequest)
					return
				}

				// Handle the event
				switch event.Type {
				case "checkout.session.completed":
					var session stripe.CheckoutSession
					err := json.Unmarshal(event.Data.Raw, &session)
					if err != nil {
						http.Error(w, "error parsing session", http.StatusBadRequest)
						return
					}

					// Payment succeeded!
					fmt.Println("Payment successful for session:", session.ID)
					// TODO: Save to database, unlock features, etc.

				case "checkout.session.async_payment_failed":
					fmt.Println("Payment failed")
					// TODO: Handle failure
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"status":"ok"}`))
				fmt.Println("hello hello hello")
			})
		})
	})

	http.ListenAndServe(":3000", r)
}
