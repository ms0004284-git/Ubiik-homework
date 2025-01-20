package main

import (
	"message-gateway-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/messages", handlers.HandleMessage)

	r.Run(":8081")
}