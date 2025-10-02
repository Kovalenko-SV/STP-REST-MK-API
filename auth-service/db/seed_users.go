package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"ksv/rest-mikroservice/auth-service/models"
	"ksv/rest-mikroservice/auth-service/utils"
)

func CreateTestUser() {

	InitDB("auth-service/data/auth.db")

	defer DB.Close()

	repo := NewUserRepository(DB)

	ctx := context.Background()

	if err := clearUsers(ctx); err != nil {
		log.Fatalf("failed to clear users: %v", err)
	}

	adminPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}

	userPassword, err := utils.HashPassword("user123")
	if err != nil {
		log.Fatalf("failed to hash user password: %v", err)
	}

	now := time.Now()

	admin := &models.User{
		ID:        uuid.NewString(),
		Login:     "admin",
		Password:  adminPassword,
		IsAdmin:   true,
		CreatedAt: now,
	}

	if err := repo.CreateUser(ctx, admin); err != nil {
		log.Fatalf("failed to create admin: %v", err)
	}

	user := &models.User{
		ID:        uuid.NewString(),
		Login:     "user",
		Password:  userPassword,
		IsAdmin:   false,
		CreatedAt: now,
	}

	if err := repo.CreateUser(ctx, user); err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	fmt.Println("Successfully cleaned table and inserted admin and user!")
}

func clearUsers(ctx context.Context) error {
	_, err := DB.ExecContext(ctx, "DELETE FROM users")
	return err
}
