package warehouse

import (
	"errors"

	"github.com/bootcamp-go/consignas-go-db.git/internal/domain"
)

var (
	ErrNotFound = errors.New("warehouse not found")
	ErrInternal = errors.New("internal error")
)

type Service interface {
	GetByID(id int) (domain.Warehouse, error)
	Create(p domain.Warehouse) (domain.Warehouse, error)
	GetAll() ([]domain.Warehouse, error)
	ReportProducts(id int) (reportProducts domain.ReportProducts, err error)
}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetByID(id int) (domain.Warehouse, error) {
	warehouse, err := s.r.GetByID(id)
	if err != nil {
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}

func (s *service) Create(p domain.Warehouse) (domain.Warehouse, error) {
	warehouse, err := s.r.Create(p)
	if err != nil {
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}

func (s *service) GetAll() ([]domain.Warehouse, error) {
	warehouses, err := s.r.GetAll()
	if err != nil {
		return []domain.Warehouse{}, err
	}
	return warehouses, nil
}

func (s *service) ReportProducts(id int) (reportProducts domain.ReportProducts, err error) {
	reportProducts, err = s.r.ReportProducts(id)
	if err != nil {
		return reportProducts, err
	}
	return reportProducts, nil
}
