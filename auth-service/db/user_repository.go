package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"ksv/rest-mikroservice/auth-service/models"
)

type UserRepository struct {
	db *DBWrapper
}

type DBWrapper struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: &DBWrapper{db},
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, login, password, is_admin, created_at, updated_at)
        VALUES (:id, :login, :password, :is_admin, :created_at, :updated_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = ?`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE login = ?`
	err := r.db.GetContext(ctx, &user, query, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.UpdatedAt = &now

	query := `
        UPDATE users SET
            login = :login,
            password = :password,
            is_admin = :is_admin,
            updated_at = :updated_at
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
