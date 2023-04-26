package handler

import (
	"errors"
	"os"
	"strconv"

	"github.com/bootcamp-go/consignas-go-db.git/internal/domain"
	"github.com/bootcamp-go/consignas-go-db.git/internal/warehouse"
	"github.com/bootcamp-go/consignas-go-db.git/pkg/web"
	"github.com/gin-gonic/gin"
)

type warehouseHandler struct {
	w warehouse.Service
}

func NewWarehouseHandler(w warehouse.Service) *warehouseHandler {
	return &warehouseHandler{
		w: w,
	}
}

func (h *warehouseHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		warehouse, err := h.w.GetByID(id)
		if err != nil {
			web.Failure(c, 404, errors.New("warehouse not found"))
			return
		}
		web.Success(c, 200, warehouse)
	}
}

func validate(warehouse *domain.Warehouse) (bool, error) {
	switch {
	case warehouse.Name == "" || warehouse.Address == "" || warehouse.Telephone == "":
		return false, errors.New("fields can't be empty")
	case warehouse.Capacity <= 0:
		return false, errors.New("capacity must be greater than 0")
	}
	return true, nil
}

func (h *warehouseHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var warehouse domain.Warehouse
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token not found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
		err := c.ShouldBindJSON(&warehouse)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
			return
		}
		valid, err := validate(&warehouse)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		w, err := h.w.Create(warehouse)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		web.Success(c, 201, w)
	}
}

func (h *warehouseHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouses, err := h.w.GetAll()
		if err != nil {
			web.Failure(c, 500, errors.New("internal error"))
			return
		}
		web.Success(c, 200, warehouses)
	}
}
func (h *warehouseHandler) ReportProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		warehouse, err := h.w.ReportProducts(id)
		if err != nil {
			web.Failure(c, 404, errors.New("warehouse not found"))
			return
		}
		web.Success(c, 200, warehouse)
	}
}
