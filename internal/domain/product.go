package domain

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	WarehouseId int     `json:"id_warehouse" binding:"required"`
}

type ProductFull struct {
	Product
	WarehouseName    string `json:"warehouse_name"`
	WarehouseAddress string `json:"warehouse_address"`
}
