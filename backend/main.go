package main

import (
	"log"
	"net/http"

	"assignment-backend/api"
	"assignment-backend/initializers"
	"assignment-backend/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Initializers
	initializers.MergeProducts()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Products
	mux.HandleFunc("/products", api.GetProducts)

	// CORS: use nil to allow all origins, or e.g. []string{"http://localhost:3000", "https://myapp.com"}
	origins := []string{} // empty = allow all
	methods := []string{http.MethodGet}
	headers := []string{"Accept", "Content-Type", "Authorization"}
	handler := middleware.CORSMiddleware(origins, methods, headers)(mux)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
