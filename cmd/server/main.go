package main

import (
	"database/sql"
	"log"

	"github.com/bootcamp-go/consignas-go-db.git/cmd/server/handler"
	"github.com/bootcamp-go/consignas-go-db.git/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {

	// storage := store.NewJsonStore("./products.json")
	// repo := product.NewRepository(storage)

	databaseConfig := mysql.Config{
		User:   "root",
		Addr:   "localhost:3306",
		DBName: "my_db",
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET(":id", productHandler.GetByID())
		products.POST("", productHandler.Post())
		products.DELETE(":id", productHandler.Delete())
		products.PATCH(":id", productHandler.Patch())
		products.PUT(":id", productHandler.Put())
	}

	r.Run(":8080")
}
