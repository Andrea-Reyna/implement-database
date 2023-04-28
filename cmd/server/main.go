package main

import (
	"database/sql"
	"log"

	"github.com/bootcamp-go/consignas-go-db.git/cmd/server/handler"
	"github.com/bootcamp-go/consignas-go-db.git/internal/product"
	"github.com/bootcamp-go/consignas-go-db.git/internal/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {

	// storage := store.NewJsonStore("./products.json")
	// repo := product.NewRepository(storage)

	databaseConfig := mysql.Config{
		User:      "root",
		Addr:      "localhost:3306",
		DBName:    "my_db",
		ParseTime: true,
	}
	database, err := sql.Open("mysql", databaseConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	defer database.Close()
	if err = database.Ping(); err != nil {
		panic(err)
	}
	log.Println("database Configured")

	repository := product.NewMySQLRepository(database)
	service := product.NewService(repository)
	productHandler := handler.NewProductHandler(service)

	warehouseRepository := warehouse.NewMySQLRepository(database)
	warehouseService := warehouse.NewService(warehouseRepository)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	products := r.Group("/products")
	{
		products.GET(":id", productHandler.GetByID())
		products.GET("", productHandler.GetAll())
		products.GET("/details/:id", productHandler.GetFullData())

		products.POST("", productHandler.Post())
		products.DELETE(":id", productHandler.Delete())
		products.PATCH(":id", productHandler.Patch())
		products.PUT(":id", productHandler.Put())
	}

	warehouses := r.Group("/warehouses")
	{
		warehouses.GET("", warehouseHandler.GetAll())
		warehouses.GET("/:id", warehouseHandler.GetByID())
		warehouses.GET("/reportProducts", warehouseHandler.ReportProducts())
		warehouses.POST("", warehouseHandler.Post())
	}

	r.Run(":8080")
}
