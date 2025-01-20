package main

import (
	"github.com/gin-gonic/gin"
	"device-service/handlers"
	"device-service/db"
	"device-service/repositories"
	"device-service/models"
)

func main() {
	client := db.InitDB()
	client.AutoMigrate(&models.Device{})

	repo := repositories.NewGormDeviceRepository(client)
	handler := handlers.NewDeviceAccess(repo)
	// deviceRepo := repositories.NewGormDeviceRepository(db.Client)

	r := gin.Default()
	r.PUT("/devices/:deviceId", handler.UpdateOrCreateDevice)
	r.GET("/devices/:deviceId/username", handler.GetDevice)

	r.Run(":8080") 
}