package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"iotservice/db"
)

type Device struct {
	ID   	 string    	`json:"id" gorm:"primaryKey  binding:"required"`
	Username string 	`json:"username"  binding:"required"`
}


func main() {
	db.InitDB()
	db.Client.AutoMigrate(&Device{}) 

	r := gin.Default()

	r.PUT("/devices/:deviceId", func(ctx *gin.Context){
		var newDevice Device
		id := ctx.Param("deviceId")
		if err := ctx.ShouldBindJSON(&newDevice); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newDevice.ID = id
		var device Device

		if err := db.Client.First(&device, "id = ?", id).Error; err == nil {
			// 存在更新
			db.Client.Model(&device).Updates(newDevice)
			ctx.JSON(http.StatusOK, gin.H{"data": device})
		} else {
			// 不存在新增
			db.Client.Create(&newDevice)
			ctx.JSON(http.StatusCreated, gin.H{"data": newDevice})
		}

	})

	r.GET("/devices/:deviceId/username", func(ctx *gin.Context){
		id := ctx.Param("deviceId")

		var device Device
		if err := db.Client.First(&device, "id = ?", id).Error; err == nil {
			ctx.JSON(http.StatusOK, gin.H{"username": device.Username})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		}
		
	})

	r.Run(":8080") 
}