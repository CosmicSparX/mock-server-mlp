package main

import (
	"fmt"
	"log"
	"net/http"

	"mock_server/internal/sites"
)

// Updated main.go

func main() {
	mux := http.NewServeMux()

	// ROUTE 1: SaaS Landing Page
	mux.HandleFunc("/", sites.SaaSHandler)

	// ROUTE 2: E-Commerce Store
	mux.HandleFunc("/shop", sites.ShopHandler)

	// ROUTE 3: News Portal
	mux.HandleFunc("/news", sites.NewsHandler)

	// DEBUG: Verification for your Proxy's Bot Detection
	mux.HandleFunc("/ua", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Your User Agent: %s", r.UserAgent())
	})

	log.Println("🚀 Mock Client Server running on http://localhost:9090")
	http.ListenAndServe(":9090", mux)
}
