package main

import (
	"RECU-CLIENTE-SERVIDOR/handler"
	"RECU-CLIENTE-SERVIDOR/repository"
	"RECU-CLIENTE-SERVIDOR/usecase"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewMemoryRepository()
	useCase := usecase.NewProductUseCase(repo)
	productHandler := handler.NewProductHandler(useCase)

	// Configurar rutas
	http.HandleFunc("/addProduct", productHandler.AddProduct)
	http.HandleFunc("/isNewProductAdded", productHandler.IsNewProductAdded)
	http.HandleFunc("/countProductsInDiscount", productHandler.CountProductsInDiscount)

	// Configurar CORS
	corsHandler := corsMiddleware(http.DefaultServeMux)

	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
