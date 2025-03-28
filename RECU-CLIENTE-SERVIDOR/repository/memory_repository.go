package repository

import (
	"RECU-CLIENTE-SERVIDOR/domain"
	"sync"
)

type MemoryRepository struct {
	products []domain.Product
	mutex    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		products: make([]domain.Product, 0),
	}
}

func (r *MemoryRepository) Save(product domain.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.products = append(r.products, product)
	return nil
}

func (r *MemoryRepository) FindRecent(since int64) ([]domain.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Solo productos creados después del timestamp
	recent := make([]domain.Product, 0)
	for _, p := range r.products {
		if p.CreatedAt > since {
			recent = append(recent, p)
		}
	}

	return recent, nil
}

func (r *MemoryRepository) CountWithDiscount() (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, p := range r.products {
		if p.Descuento {
			count++
		}
	}

	return count, nil
}
