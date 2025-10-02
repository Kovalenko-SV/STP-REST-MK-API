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

	"ksv/rest-mikroservice/auth-service/models"
)

func setupTestDB(t *testing.T) *sqlx.DB {

	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)

	schema := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	);`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Створення нового користувача
	user := &models.User{
		ID:        uuid.NewString(),
		Login:     "testuser",
		Password:  "testpassword",
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}

	t.Run("CreateUser", func(t *testing.T) {
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		fetchedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Login, fetchedUser.Login)
	})

	t.Run("GetUserByLogin", func(t *testing.T) {
		fetchedUser, err := repo.GetUserByLogin(ctx, user.Login)
		require.NoError(t, err)
		assert.Equal(t, user.ID, fetchedUser.ID)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user.Login = "updateduser"
		user.Password = "newpassword"
		err := repo.UpdateUser(ctx, user)
		require.NoError(t, err)

		updatedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "updateduser", updatedUser.Login)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := repo.DeleteUser(ctx, user.ID)
		require.NoError(t, err)

		_, err = repo.GetUserByID(ctx, user.ID)
		require.Error(t, err)
	})
}
