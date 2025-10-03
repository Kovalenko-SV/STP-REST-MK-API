package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ksv/rest-mikroservice/auth-service/db"
	"ksv/rest-mikroservice/auth-service/models"
	"ksv/rest-mikroservice/auth-service/token"
	"ksv/rest-mikroservice/auth-service/utils"

	"github.com/joho/godotenv"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Login   string `json:"login"`
	IsAdmin bool   `json:"is_admin"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// NewValidateTokenHandler godoc
// @Summary Валідація токена
// @Description Перевіряє дійсність JWT токена
// @Tags auth
// @Produce plain
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Invalid token"
// @Security BearerAuth
// @Router /auth/validate [get]
func NewValidateTokenHandler(secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenStr := headerParts[1]

		jwtMaker := token.NewJWTMaker(secretKey)

		_, err := jwtMaker.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// NewAuthHandler godoc
// @Summary Авторизація користувача
// @Description Приймає логін і пароль, перевіряє користувача та повертає JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Дані користувача"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} models.ErrorResponse "Некоректний запит"
// @Failure 401 {object} models.ErrorResponse "Невірні дані"
// @Failure 500 {object} models.ErrorResponse "Помилка сервера"
// @Router /api/auth [post]
func NewAuthHandler(secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request"})
			return
		}

		repo := db.NewUserRepository(db.DB)
		user, err := repo.GetUserByLogin(context.Background(), req.Login)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Database error"})
			return
		}

		if err := utils.CheckPassword(req.Password, user.Password); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid credentials"})
			return
		}

		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found, relying on system environment variables")
		}
		tokenDurationStr := os.Getenv("JWT_EXPIRET_TIME")
		if tokenDurationStr == "" {
			log.Println("Warning: .env ket – JWT_EXPIRET_TIME not found")
		}
		tokenDuration, err := time.ParseDuration(tokenDurationStr)
		if err != nil {
			log.Println("Time durattion set error")
		}

		jwtMaker := token.NewJWTMaker(secretKey)

		tokenString, claims, err := jwtMaker.CreateToken(user.ID, user.Login, user.IsAdmin, tokenDuration)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to create token"})
			return
		}

		userResponse := UserResponse{
			ID:      user.ID,
			Login:   user.Login,
			IsAdmin: user.IsAdmin,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{
			Token:     tokenString,
			ExpiresAt: claims.ExpiresAt.Time,
			User:      userResponse,
		})
	}
}
