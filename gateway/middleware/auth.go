package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	_ "ksv/rest-mikroservice/auth-service/models"
)

const (
	authServiceURL = "http://localhost:8081"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			publicPaths := map[string]bool{
				"/api/auth":           true,
				"/swagger/":           true,
				"/swagger/index.html": true,
			}

			for pub := range publicPaths {
				if strings.HasPrefix(r.URL.Path, pub) {
					log.Println("Public path, no token required.")
					next.ServeHTTP(w, r)
					return
				}
			}

			if publicPaths[r.URL.Path] {
				log.Println("Public path, no token required.")
				next.ServeHTTP(w, r)
				log.Println(http.StatusOK, r.Method, r.URL.Path, time.Since(start))
				return
			}

			req, err := http.NewRequest(http.MethodGet, authServiceURL+"/auth/validate", nil)
			if err != nil {
				http.Error(w, "Could not create auth request", http.StatusInternalServerError)
				return
			}

			req.Header.Set("Authorization", r.Header.Get("Authorization"))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "Error contacting auth service", http.StatusInternalServerError)
				log.Println("Auth service error:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Println(http.StatusUnauthorized, r.Method, r.URL.Path, time.Since(start))
				return
			}

			ctx := context.Background()
			r = r.WithContext(ctx)

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r.WithContext(ctx))
			log.Println(http.StatusOK, r.Method, r.URL.Path, time.Since(start))
		})
	}
}
