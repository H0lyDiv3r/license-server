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

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	fmt.Println("webook", webhookSecret)

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
	r.Post("/stripeWebhook", func(w http.ResponseWriter, r *http.Request) {

		signature := r.Header.Get("Stripe-Signature")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			http.Error(w, "error reading body", http.StatusBadRequest)
			return
		}

		fmt.Println("Signature:", signature)
		fmt.Println("Body length:", len(body))

		webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		fmt.Println("Webhook secret loaded:", webhookSecret != "")

		event, err := webhook.ConstructEventWithOptions(body, signature, webhookSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
		if err != nil {
			fmt.Println("Signature verification failed:", err)
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		fmt.Println("Event type:", event.Type)

		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &session)
			if err != nil {
				fmt.Println("Error parsing session:", err)
				http.Error(w, "error parsing session", http.StatusBadRequest)
				return
			}
			fmt.Println("✓ Payment successful for session:", session.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
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
		})
	})

	http.ListenAndServe(":3000", r)
}
