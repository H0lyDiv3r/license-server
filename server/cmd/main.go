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
		sessionID := r.URL.Query().Get("session_id")

		html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Payment Successful</title>
  <style>
    * { box-sizing: border-box; }
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
    h1 { margin: 0 0 12px; font-size: 2rem; font-weight: 700; }
    p { margin: 0; color: #a3a3a3; line-height: 1.6; }
    .spinner {
      margin: 32px auto;
      width: 40px;
      height: 40px;
      border: 4px solid #2a2a2a;
      border-top: 4px solid #f5f5f5;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }
    @keyframes spin { to { transform: rotate(360deg); } }
    .key-box {
      margin-top: 22px;
      padding: 16px;
      background: #2a2a2a;
      border-radius: 8px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      font-size: 1rem;
      color: #e5e5e5;
      word-break: break-all;
      user-select: all;
    }
    .copy-btn {
      margin-top: 16px;
      height: 38px;
      padding: 0 24px;
      border: 0;
      border-radius: 8px;
      background: #f5f5f5;
      color: #111111;
      font-size: 0.9rem;
      font-weight: 600;
      cursor: pointer;
    }
    .copy-btn:hover { background: #e5e5e5; }
    .status { margin-top: 12px; min-height: 18px; font-size: 0.9rem; color: #a3a3a3; }
    .hidden { display: none; }
  </style>
</head>
<body>
  <div class="container">
    <div id="loading">
      <h1>Processing Payment</h1>
      <p>Getting your license key...</p>
      <div class="spinner"></div>
    </div>

    <div id="content" class="hidden">
      <h1>Payment Successful</h1>
      <p>Your license key is below. Copy it and paste it into the app.</p>
      <div class="key-box" id="licenseKey"></div>
      <button class="copy-btn" id="copyBtn">Copy to Clipboard</button>
      <div class="status" id="status"></div>
    </div>

    <div id="error" class="hidden">
      <h1>Something went wrong</h1>
      <p>Could not find your license key. Please contact support.</p>
    </div>
  </div>

  <script>
    const sessionId = '` + sessionID + `';

    if (!sessionId) {
      document.getElementById("loading").classList.add("hidden");
      document.getElementById("error").classList.remove("hidden");
    } else {
      let attempts = 0;
      const maxAttempts = 150; // 30 seconds

      const interval = setInterval(() => {
        attempts++;
        if (attempts > maxAttempts) {
          clearInterval(interval);
          document.getElementById("loading").classList.add("hidden");
          document.getElementById("error").classList.remove("hidden");
          return;
        }

        fetch("/licenseBySession/" + sessionId)
          .then((res) => res.json())
          .then((data) => {
            if (data && data.key) {
              clearInterval(interval);
              document.getElementById("loading").classList.add("hidden");
              document.getElementById("content").classList.remove("hidden");
              document.getElementById("licenseKey").textContent = data.key;
            }
          })
          .catch(() => {});
      }, 200);

      document.getElementById("copyBtn").addEventListener("click", async () => {
        const key = document.getElementById("licenseKey").textContent;
        try {
          await navigator.clipboard.writeText(key);
          document.getElementById("status").textContent = "Copied!";
        } catch {
          document.getElementById("status").textContent = "Copy failed.";
        }
      });
    }
  </script>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})
	// r.Get("/failure")

	r.Get("/licenseBySession/{sessionID}", func(w http.ResponseWriter, r *http.Request) {
		sessionID := chi.URLParam(r, "sessionID")
		license, err := licenseService.GetLicenseBySessionId(r.Context(), sessionID)
		if err != nil {
			http.Error(w, "license not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(license)
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
