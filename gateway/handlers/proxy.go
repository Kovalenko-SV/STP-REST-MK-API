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

func ProxyAuthService(w http.ResponseWriter, r *http.Request) {
	proxyRequest(authServiceURL, w, r)
}

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
