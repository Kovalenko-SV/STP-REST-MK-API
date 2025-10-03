package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	authServiceURL    = "http://localhost:8081"
	productServiceURL = "http://localhost:8082"
)

// ProxyAuthService godoc
// @Summary Авторизація користувача
// @Description Приймає логін і пароль, перевіряє користувача та повертає JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body handlers.AuthRequest true "Дані користувача"
// @Success 200 {object} handlers.AuthResponse
// @Failure 400 {object} handlers.ErrorResponse "Некоректний запит"
// @Failure 401 {object} handlers.ErrorResponse "Невірні дані"
// @Failure 500 {object} handlers.ErrorResponse "Помилка сервера"
// @Router /api/auth [post]
func ProxyAuthService(w http.ResponseWriter, r *http.Request) {
	proxyRequest(authServiceURL, w, r)
}

// ProxyProductService godoc
// @Summary Операції з продуктами (проксі)
// @Description Проксі-ендпоінт, який передає запити у product-service.
// @Tags product
// @Accept json
// @Produce json
//
// @Param id query string false "ID продукту (для GET, PUT, DELETE)"
// @Param limit query int false "Кількість продуктів (для GET, за замовчуванням 10)"
//
// @Security BearerAuth
//
// @Success 200 {object} handlers.Product
// @Success 200 {array} handlers.Product
// @Success 201 {object} handlers.Product
// @Success 200 {object} handlers.Product
//
// @Failure 400 {string} string "Некоректний запит"
// @Failure 404 {string} string "Продукт не знайдено"
// @Failure 500 {string} string "Помилка сервера"
//
// @Security BearerAuth
// @Router /api/product [get]
// @Router /api/product [delete]
//
// @Router /api/product [post]  [x-product-body handlers.Product]
// @Router /api/product [put]   [x-product-body handlers.Product]
func ProxyProductService(w http.ResponseWriter, r *http.Request) {
	proxyRequest(productServiceURL, w, r)
}

func proxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing target URL: %v", err), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = target.Host

	proxy.ServeHTTP(w, r)
}
