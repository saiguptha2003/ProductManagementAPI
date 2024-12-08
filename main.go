package main

import (
	"ProductManagement/cache"
	"ProductManagement/db"
	"ProductManagement/handlers"
	"ProductManagement/queue"
	"ProductManagement/services"
	"log"

	"github.com/gin-gonic/gin"
)


func main() {
	db.InitDB()
	cache.InitRedis()
	queue.InitQueue()
	go services.StartImageProcessor(queue.Channel)
	r := gin.Default()
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProductByID)
	r.GET("/products", handlers.GetProductsHandler)
	log.Fatal(r.Run(":8080"))
};