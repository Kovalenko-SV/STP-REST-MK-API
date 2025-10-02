package db

import (
	"context"
	"fmt"
	"ksv/rest-mikroservice/product-service/models"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CreateTestProducts() {

	InitDB("./data/products.db")
	defer DB.Close()

	repo := NewProductRepository(DB)
	ctx := context.Background()

	if err := clearProducts(ctx); err != nil {
		log.Fatalf("failed to clear products: %v", err)
	}

	now := time.Now()

	product1 := &models.Product{
		ID:        uuid.NewString(),
		Name:      "Laptop",
		Price:     1200.50,
		Quantity:  10,
		CreatedAt: now,
	}

	if err := repo.CreateProduct(ctx, product1); err != nil {
		log.Fatalf("failed to create product1: %v", err)
	}

	product2 := &models.Product{
		ID:        uuid.NewString(),
		Name:      "Phone",
		Price:     700.00,
		Quantity:  25,
		CreatedAt: now,
	}

	if err := repo.CreateProduct(ctx, product2); err != nil {
		log.Fatalf("failed to create product2: %v", err)
	}

	fmt.Println("Successfully cleaned table and inserted test products!")
}

func clearProducts(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, "DELETE FROM products")
	return err
}
