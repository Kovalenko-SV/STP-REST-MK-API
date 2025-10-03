// @title Auth Service API
// @version 1.0
// @description API для авторизації користувачів
// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"

	_ "ksv/rest-mikroservice/auth-service/docs"

	"github.com/ianschenck/envflag"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/swaggo/files"

	"ksv/rest-mikroservice/auth-service/db"
	"ksv/rest-mikroservice/auth-service/handlers"
)

const (
	port = ":8081"
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

	db.InitDB("./data/auth.db")

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

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /api/auth", handlers.NewAuthHandler(secretKey))
	mux.HandleFunc("GET /auth/validate", handlers.NewValidateTokenHandler(secretKey))

	return mux
}
