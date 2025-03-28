package handler

import (
	"RECU-CLIENTE-SERVIDOR/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ProductHandler struct {
	useCase             *usecase.ProductUseCase
	discountSubscribers []chan int
	mutex               sync.Mutex
}

func NewProductHandler(useCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		useCase:             useCase,
		discountSubscribers: make([]chan int, 0),
	}
}

type AddProductRequest struct {
	Nombre    string `json:"nombre"`
	Precio    int    `json:"precio"`
	Codigo    string `json:"codigo"`
	Descuento bool   `json:"descuento"`
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Recibido nuevo producto: %+v\n", r)

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req AddProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	product, err := h.useCase.AddProduct(req.Nombre, req.Precio, req.Codigo, req.Descuento)
	if err != nil {
		http.Error(w, "Error al agregar producto", http.StatusInternalServerError)
		return
	}

	if req.Descuento {
		h.notifyDiscountSubscribers()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) IsNewProductAdded(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	sinceStr := r.URL.Query().Get("since")
	var since int64 = 0
	if sinceStr != "" {
		since, _ = strconv.ParseInt(sinceStr, 10, 64)
	} else {
		since = time.Now().Unix() - 30
	}

	products, err := h.useCase.GetRecentProducts(since)
	if err != nil {
		http.Error(w, "Error al obtener productos recientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"now":      time.Now().Unix(),
	})
}

func (h *ProductHandler) CountProductsInDiscount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	count, err := h.useCase.GetDiscountedProductsCount()
	if err != nil {
		http.Error(w, "Error al contar productos con descuento", http.StatusInternalServerError)
		return
	}

	expectedCountStr := r.URL.Query().Get("expectedCount")
	if expectedCountStr != "" {
		expectedCount, err := strconv.Atoi(expectedCountStr)
		if err == nil && expectedCount == count {
			// Long polling
			ch := make(chan int, 1)

			h.mutex.Lock()
			h.discountSubscribers = append(h.discountSubscribers, ch)
			h.mutex.Unlock()

			timeout := time.After(30 * time.Second)

			select {
			case newCount := <-ch:
				json.NewEncoder(w).Encode(map[string]int{"count": newCount})
			case <-timeout:
				json.NewEncoder(w).Encode(map[string]int{"count": count})
			}

			close(ch)
			h.removeSubscriber(ch)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func (h *ProductHandler) notifyDiscountSubscribers() {
	count, err := h.useCase.GetDiscountedProductsCount()
	if err != nil {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, ch := range h.discountSubscribers {
		select {
		case ch <- count:
			// Enviado correctamente
		default:
			// Canal lleno o cerrado, ignorar
		}
	}
}

func (h *ProductHandler) removeSubscriber(ch chan int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i, subscriber := range h.discountSubscribers {
		if subscriber == ch {
			h.discountSubscribers = append(h.discountSubscribers[:i], h.discountSubscribers[i+1:]...)
			break
		}
	}
}
