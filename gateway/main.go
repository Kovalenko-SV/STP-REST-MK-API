// @title Gateway API
// @version 1.0
// @description Проксі API для auth та product сервісів
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"
	"net/http"

	"ksv/rest-mikroservice/gateway/handlers"
	"ksv/rest-mikroservice/gateway/middleware"

	_ "ksv/rest-mikroservice/gateway/docs"

	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	port = ":8080"
)

func main() {

	router := setupRouter()

	server := http.Server{
		Addr:    port,
		Handler: middleware.AuthMiddleware()(router),
	}

	fmt.Printf("Gateway starting on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}

func setupRouter() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/api/auth", handlers.ProxyAuthService)

	mux.HandleFunc("/api/product", handlers.ProxyProductService)

	return mux
}
