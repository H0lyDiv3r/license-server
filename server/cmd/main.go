package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"license-server/db"
	"license-server/internals/domain"
	"license-server/internals/handler"
	"license-server/internals/middleware"
	"license-server/internals/repository"
	"license-server/internals/service"
	"license-server/payment"
	"log"
	"net/http"
	"os"
	"strconv"

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

	paymentRepo := repository.NewpaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
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
			userId := session.Metadata["user_id"]
			email := session.Metadata["user_email"]
			duration, err := strconv.Atoi(session.Metadata["duration"])
			if err != nil {
				http.Error(w, "error cant read duration", http.StatusBadRequest)
				return
			}

			uid, err := strconv.Atoi(userId)
			if err != nil {
				http.Error(w, "error cant read duration", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "user", domain.UserPayload{
				UserId: float64(uid),
				Email:  email,
			})
			// now that payment succeeded.
			// generate the license key and store it
			license, err := licenseService.GenerateKey(ctx, domain.GenerateKeyRequest{
				Duration: uint(duration),
			})
			if err != nil {
				http.Error(w, "error cant generate key"+err.Error(), http.StatusBadRequest)
				return
			}

			// need to store the payment session
			payment := domain.SavePaymentRequest{
				SessionId:     session.ID,
				Amount:        uint(session.AmountSubtotal),
				Currency:      string(session.Currency),
				Email:         session.CustomerEmail,
				Status:        string(session.Status),
				PaymentMethod: "Card",
				Metadata:      "",
				LicenseId:     license.ID,
				// UserId:        uint(uid),
			}

			pSes, err := paymentService.SavePaymentSession(r.Context(), payment)
			if err != nil {
				http.Error(w, "error cant generate key"+err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println("✓ Payment successful for session:", session, userId, email, pSes)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/success", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("bad stuff happens in the bathroom", r.URL.Query())
		// Get the session ID from the URL parameter
		// TODO: Look up the session in your database to get the license key
		// For now, just show a placeholder

		html := `
		<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Payment Successful</title>
  <style>
    * {
      box-sizing: border-box;
    }

    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      padding: 24px;
      background: #111111;
      color: #f5f5f5;
      font-family: Arial, sans-serif;
    }

    .container {
      width: 100%;
      max-width: 560px;
      text-align: center;
      padding: 24px;
    }

    h1 {
      margin: 0 0 12px;
      font-size: 2rem;
      font-weight: 700;
    }

    p {
      margin: 0;
      color: #a3a3a3;
      line-height: 1.6;
    }

    .key-row {
      margin-top: 22px;
      display: flex;
      align-items: center;
      gap: 10px;
      justify-content: center;
    }

    .key {
      flex: 1;
      min-width: 0;
      padding: 12px 0;
      background: transparent;
      border: none;
      color: #e5e5e5;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      font-size: 0.95rem;
      text-align: left;
      word-break: break-all;
      user-select: all;
    }

    button {
      flex: 0 0 auto;
      height: 34px;
      padding: 0 12px;
      border: 0;
      border-radius: 8px;
      background: #2a2a2a;
      color: #f5f5f5;
      font-size: 0.85rem;
      cursor: pointer;
    }

    button:hover {
      background: #3a3a3a;
    }

    .status {
      margin-top: 12px;
      min-height: 18px;
      font-size: 0.9rem;
      color: #b5b5b5;
    }

    @media (max-width: 480px) {
      .key-row {
        flex-direction: column;
        align-items: stretch;
      }

      .key {
        text-align: center;
        padding: 8px 0;
      }

      button {
        width: 100%;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Payment Successful</h1>
    <p>Your license key is below.</p>

    <div class="key-row">
      <div class="key" id="licenseKey">` + key + r.URL.RawQuery + `</div>
      <button id="copyBtn">Copy</button>
    </div>

    <div class="status" id="status"></div>
  </div>

  <script>
    const copyBtn = document.getElementById("copyBtn");
    const licenseKey = document.getElementById("licenseKey");
    const status = document.getElementById("status");

    copyBtn.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText(licenseKey.textContent.trim());
        status.textContent = "Copied.";
      } catch {
        status.textContent = "Copy failed.";
      }
    });
  </script>
</body>
</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
	// r.Get("/failure")

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
