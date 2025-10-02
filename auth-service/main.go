package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ianschenck/envflag"
	"github.com/joho/godotenv"

	"ksv/rest-mikroservice/auth-service/db"
	"ksv/rest-mikroservice/auth-service/handlers"
)

const (
	port = "0.0.0.0:8081"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	var secretKey = envflag.String("JWT_SECRET_KEY", "01234567890123456789012345678901", "Cекретний ключ, що використовується для підписання JWT ")
	if len(*secretKey) < 32 {
		log.Fatal("JWT_SECRET_KEY must be at least 32 characters long")
	}

	envflag.Parse()

	db.InitDB("auth-service/data/auth.db")

	defer db.DB.Close()

	router := setupRouter(*secretKey)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	envflag.Parse()

	fmt.Printf("Auth service starting on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}

func setupRouter(secretKey string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", handlers.NewAuthHandler(secretKey))
	mux.HandleFunc("GET /auth/validate", handlers.NewValidateTokenHandler(secretKey))

	return mux
}
