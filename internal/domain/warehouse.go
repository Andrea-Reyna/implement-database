package domain

type Warehouse struct {
	Id        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Telephone string `json:"telephone" binding:"required"`
	Capacity  int    `json:"capacity"`
}
