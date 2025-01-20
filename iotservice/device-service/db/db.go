package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Client *gorm.DB

func InitDB() *gorm.DB {
	dsn := "root:root@tcp(mysql:3306)/iot?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Client, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return Client
}

