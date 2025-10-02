package main

import (
	"fmt"
	"log"
	"net/http"

	"ksv/rest-mikroservice/gateway/handlers"
	"ksv/rest-mikroservice/gateway/middleware"
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

	mux.HandleFunc("/api/auth", handlers.ProxyAuthService)

	mux.HandleFunc("/api/product", handlers.ProxyProductService)

	return mux
}
