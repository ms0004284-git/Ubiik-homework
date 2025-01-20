package repositories

import (
	"errors"
	"gorm.io/gorm"
	"device-service/models"
)

type DeviceRepository interface {
	GetDeviceByID(id string) (*models.Device, error)
	CreateDevice(device *models.Device) error
	UpdateDevice(device *models.Device) error
}

type GormDeviceRepository struct {
	DB *gorm.DB
}

func NewGormDeviceRepository(db *gorm.DB) *GormDeviceRepository {
	return &GormDeviceRepository{DB: db}
}

func (repo *GormDeviceRepository) GetDeviceByID(id string) (*models.Device, error) {
	var device models.Device
	if err := repo.DB.First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (repo *GormDeviceRepository) CreateDevice(device *models.Device) error {
	return repo.DB.Create(device).Error
}

func (repo *GormDeviceRepository) UpdateDevice(device *models.Device) error {
	result := repo.DB.Model(&models.Device{}).
		Where("id = ? AND version = ?", device.ID, device.Version).
		Updates(map[string]interface{}{
			"username": device.Username,
			"version":  device.Version + 1, // optimistic locking, 更新時增加 version 
		})
	if result.RowsAffected == 0 {
		return errors.New("update failed: version mismatch or device not found")
	}

	return result.Error
}
