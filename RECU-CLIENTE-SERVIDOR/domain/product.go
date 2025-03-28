package domain

type Product struct {
	Nombre    string
	Precio    int
	Codigo    string
	Descuento bool
	CreatedAt int64
}

type ProductRepository interface {
	Save(product Product) error
	FindRecent(since int64) ([]Product, error)
	CountWithDiscount() (int, error)
}
