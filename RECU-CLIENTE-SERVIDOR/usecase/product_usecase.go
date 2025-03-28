package usecase

import (
	"RECU-CLIENTE-SERVIDOR/domain"
	"time"
)

type ProductUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) AddProduct(nombre string, precio int, codigo string, descuento bool) (domain.Product, error) {
	product := domain.Product{
		Nombre:    nombre,
		Precio:    precio,
		Codigo:    codigo,
		Descuento: descuento,
		CreatedAt: time.Now().Unix(),
	}

	err := uc.repo.Save(product)
	return product, err
}

func (uc *ProductUseCase) GetRecentProducts(since int64) ([]domain.Product, error) {
	return uc.repo.FindRecent(since)
}

func (uc *ProductUseCase) GetDiscountedProductsCount() (int, error) {
	return uc.repo.CountWithDiscount()
}
