package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"device-service/repositories"
	"device-service/models"
	"fmt"
)

type DeviceAccess struct {
	Repo repositories.DeviceRepository
}

func NewDeviceAccess(repo repositories.DeviceRepository) *DeviceAccess {
	return &DeviceAccess{Repo: repo}
}

func (access *DeviceAccess) UpdateOrCreateDevice(ctx *gin.Context) {
	var newDevice models.Device
	id := ctx.Param("deviceId")
	if err := ctx.ShouldBindJSON(&newDevice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDevice.ID = id
	device, err := access.Repo.GetDeviceByID(id)
	fmt.Print(err)
	if err == nil {
		// 更新
		device.Username = newDevice.Username
		if err := access.Repo.UpdateDevice(device); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Device updating error."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Device username updated."})
		return
	}
	
	// 新增
	if err := access.Repo.CreateDevice(&newDevice); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Device creating error."})
		return
	}
	
	ctx.JSON(http.StatusCreated, gin.H{"message": "Device created"})

}

func (access *DeviceAccess) GetDevice(ctx *gin.Context) {
	id := ctx.Param("deviceId")
	device, err := access.Repo.GetDeviceByID(id)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"username": device.Username})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Item not found."})
	}
}