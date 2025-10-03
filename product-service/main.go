// @title Product Service API
// @version 1.0
// @description API для роботи з продуктами
// @host localhost:8082
// @BasePath /
package main

import (
	"fmt"
	"log"
	"net/http"

	"ksv/rest-mikroservice/product-service/db"
	"ksv/rest-mikroservice/product-service/handlers"

	_ "ksv/rest-mikroservice/product-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/swaggo/files"
)

const port = ":8082"

func main() {
	db.InitDB("./data/products.db")
	defer db.DB.Close()

	repo := db.NewProductRepository(db.DB)
	// db.CreateTestProducts()
	productHandler := handlers.NewProductHandler(repo)

	router := setupRouter(productHandler)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("Product service starting on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(h *handlers.ProductHandler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /api/product", h.Create())
	mux.HandleFunc("GET /api/product", h.Get())
	mux.HandleFunc("PUT /api/product", h.Update())
	mux.HandleFunc("DELETE /api/product", h.Delete())
	return mux
}
