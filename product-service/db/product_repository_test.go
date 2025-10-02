package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ksv/rest-mikroservice/product-service/models"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)

	schema := `
	CREATE TABLE products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	);`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestProductRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	product := &models.Product{
		ID:        uuid.NewString(),
		Name:      "Test Product",
		Price:     99.99,
		Quantity:  5,
		CreatedAt: time.Now(),
	}

	t.Run("CreateProduct", func(t *testing.T) {
		err := repo.CreateProduct(ctx, product)
		require.NoError(t, err)
	})

	t.Run("GetProductByID", func(t *testing.T) {
		fetchedProduct, err := repo.GetProductByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, product.Name, fetchedProduct.Name)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		product.Name = "Updated Product"
		product.Price = 150.00
		product.Quantity = 20
		err := repo.UpdateProduct(ctx, product)
		require.NoError(t, err)

		updatedProduct, err := repo.GetProductByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Product", updatedProduct.Name)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		err := repo.DeleteProduct(ctx, product.ID)
		require.NoError(t, err)

		_, err = repo.GetProductByID(ctx, product.ID)
		require.Error(t, err)
	})
}
