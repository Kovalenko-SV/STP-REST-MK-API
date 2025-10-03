package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"ksv/rest-mikroservice/product-service/db"
	"ksv/rest-mikroservice/product-service/models"

	"github.com/google/uuid"
)

type ProductHandler struct {
	repo *db.ProductRepository
}

func NewProductHandler(repo *db.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Create godoc
// @Summary Створення нового продукту
// @Description Додає новий продукт у базу даних
// @Tags product
// @Accept json
// @Produce json
// @Param product body models.Product true "Новий продукт"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Некоректний запит"
// @Failure 500 {string} string "Помилка сервера"
// @Router /api/product [post]
func (h *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.ID = uuid.NewString()
		p.CreatedAt = time.Now()
		if err := h.repo.CreateProduct(context.Background(), &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

// Get godoc
// @Summary Отримання продукту або списку продуктів
// @Description Якщо передано id — повертає один продукт, інакше список з обмеженням limit
// @Tags product
// @Produce json
// @Param id query string false "ID продукту"
// @Param limit query int false "Кількість продуктів (за замовчуванням 10)"
// @Success 200 {object} models.Product
// @Success 200 {array} models.Product
// @Failure 404 {string} string "Продукт не знайдено"
// @Failure 500 {string} string "Помилка сервера"
// @Router /api/product [get]
func (h *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id != "" {
			product, err := h.repo.GetProductByID(context.Background(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(product)
			return
		}
		limit := 10
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil {
				limit = v
			}
		}
		products, err := h.repo.GetProducts(context.Background(), limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)
	}
}

// Update godoc
// @Summary Оновлення продукту
// @Description Оновлює дані продукту за ID
// @Tags product
// @Accept json
// @Produce json
// @Param id query string true "ID продукту"
// @Param product body models.Product true "Оновлені дані продукту"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {string} string "Некоректний запит"
// @Failure 500 {string} string "Помилка сервера"
// @Router /api/product [put]
func (h *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.ID = id
		if err := h.repo.UpdateProduct(context.Background(), &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Product updated"})
	}
}

// Delete godoc
// @Summary Видалення продукту
// @Description Видаляє продукт за ID
// @Tags product
// @Produce json
// @Param id query string true "ID продукту"
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {string} string "Помилка сервера"
// @Router /api/product [delete]
func (h *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if err := h.repo.DeleteProduct(context.Background(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Product deleted"})
	}
}
