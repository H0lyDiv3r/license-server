package main

import (
	"license-server/db"
	"license-server/internals/handler"
	"license-server/internals/middleware"
	"license-server/internals/repository"
	"license-server/internals/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {

	db := db.Connect()
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewHandler(userService)

	licenseRepo := repository.NewLicenseRepository(db)
	licenseService := service.NewLicenseService(licenseRepo)
	licenseHandler := handler.NewLicenseHandler(licenseService, userService)

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
	})

	http.ListenAndServe(":3000", r)
}
