package main

import (
	"fmt"
	"license-server/db"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {

	db.Connect()
	r := chi.NewRouter()

	usr := os.Getenv("POSTGRES_USER")
	fmt.Println("main working", usr)
	http.ListenAndServe(":3000", r)
}
