package product

import (
	"errors"

	"github.com/bootcamp-go/consignas-go-db.git/internal/domain"
	"github.com/bootcamp-go/consignas-go-db.git/pkg/store"
)

type Repository interface {
	// GetByID busca un producto por su id
	GetByID(id int) (domain.Product, error)
	// GetAll busca todos los productos
	GetAll() ([]domain.Product, error)
	// GetAll busca todos los productos y agrega datos de warehouse
	GetFullData(id int) (domain.ProductFull, error)
	// Create agrega un nuevo producto
	Create(p domain.Product) (domain.Product, error)
	// Update actualiza un producto
	Update(id int, p domain.Product) (domain.Product, error)
	// Delete elimina un producto
	Delete(id int) error
}

type repository struct {
	storage store.StoreInterface
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.StoreInterface) Repository {
	return &repository{storage}
}

func (r *repository) GetByID(id int) (domain.Product, error) {
	product, err := r.storage.Read(id)
	if err != nil {
		return domain.Product{}, errors.New("product not found")
	}
	return product, nil

}

// Only implemented in mysql_repository
func (r *repository) GetAll() ([]domain.Product, error) {
	return []domain.Product{}, nil
}

// Only implemented in mysql_repository
func (r *repository) GetFullData(id int) (domain.ProductFull, error) {
	return domain.ProductFull{}, nil
}

func (r *repository) Create(p domain.Product) (domain.Product, error) {
	if !r.storage.Exists(p.CodeValue) {
		return domain.Product{}, errors.New("code value already exists")
	}
	err := r.storage.Create(p)
	if err != nil {
		return domain.Product{}, errors.New("error creating product")
	}
	return p, nil
}

func (r *repository) Delete(id int) error {
	err := r.storage.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(id int, p domain.Product) (domain.Product, error) {
	if !r.storage.Exists(p.CodeValue) {
		return domain.Product{}, errors.New("code value already exists")
	}
	err := r.storage.Update(p)
	if err != nil {
		return domain.Product{}, errors.New("error updating product")
	}
	return p, nil
}
