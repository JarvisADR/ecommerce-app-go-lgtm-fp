package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
	})

	// Product Service Reverse Proxy
	productURL, err := url.Parse("http://product-service.ecommerce.svc.cluster.local:8080")
	if err != nil {
		log.Fatal("Error parsing product service URL:", err)
	}
	productProxy := httputil.NewSingleHostReverseProxy(productURL)
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/products"
		r.RequestURI = ""
		productProxy.ServeHTTP(w, r)
	})

	// Order Service Reverse Proxy
	orderURL, err := url.Parse("http://order-service.ecommerce.svc.cluster.local:8080")
	if err != nil {
		log.Fatal("Error parsing order service URL:", err)
	}
	orderProxy := httputil.NewSingleHostReverseProxy(orderURL)
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/orders"
		r.RequestURI = ""
		orderProxy.ServeHTTP(w, r)
	})

	// User Service Reverse Proxy
	userURL, err := url.Parse("http://user-service.ecommerce.svc.cluster.local:8080")
	if err != nil {
		log.Fatal("Error parsing user service URL:", err)
	}
	userProxy := httputil.NewSingleHostReverseProxy(userURL)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/users"
		r.RequestURI = ""
		userProxy.ServeHTTP(w, r)
	})

	// Payment Service Reverse Proxy
	paymentURL, err := url.Parse("http://payment-service.ecommerce.svc.cluster.local:8080")
	if err != nil {
		log.Fatal("Error parsing payment service URL:", err)
	}
	paymentProxy := httputil.NewSingleHostReverseProxy(paymentURL)
	mux.HandleFunc("/api/payments", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/payments"
		r.RequestURI = ""
		paymentProxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}