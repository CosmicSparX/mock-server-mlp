package main

import (
	"fmt"
	"log"
	"net/http"

	"mock_server/internal/db"
	"mock_server/internal/sites"
)

// Updated main.go

func main() {
	// Initialize SQLite Database
	db.InitDB()
	defer db.DB.Close()

	mux := http.NewServeMux()

	// ROUTE 0: Minimalist Directory Home
	mux.HandleFunc("/", sites.HomeHandler)

	// ROUTE 1: SaaS Landing Page (and subroutes like /login, /features)
	mux.HandleFunc("/saas/", sites.SaaSHandler)
	mux.HandleFunc("/saas", sites.SaaSHandler)

	// ROUTE 2: E-Commerce Store
	mux.HandleFunc("/shop/", sites.ShopHandler)
	mux.HandleFunc("/shop", sites.ShopHandler)

	// ROUTE 3: News Portal
	mux.HandleFunc("/news/", sites.NewsHandler)
	mux.HandleFunc("/news", sites.NewsHandler)

	// DEBUG: Verification for your Proxy's Bot Detection
	mux.HandleFunc("/ua", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Your User Agent: %s", r.UserAgent())
	})

	log.Println("🚀 Mock Client Server running on http://localhost:9090")
	http.ListenAndServe(":9090", mux)
}
