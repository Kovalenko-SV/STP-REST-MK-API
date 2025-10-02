package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"ksv/rest-mikroservice/product-service/models"
)

type ProductRepository struct {
	db *DBWrapper
}

type DBWrapper struct {
	*sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: &DBWrapper{db},
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
        INSERT INTO products (id, name, price, quantity, created_at, updated_at)
        VALUES (:id, :name, :price, :quantity, :created_at, :updated_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	query := `SELECT * FROM products WHERE id = ?`
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	query := `
        SELECT * FROM products 
        ORDER BY created_at DESC 
        LIMIT ?
    `
	err := r.db.SelectContext(ctx, &products, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	now := time.Now()
	product.UpdatedAt = &now

	query := `
        UPDATE products SET
            name = :name,
            price = :price,
            quantity = :quantity,
            updated_at = :updated_at
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
