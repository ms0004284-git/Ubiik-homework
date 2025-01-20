package models

type Device struct {
	ID   	 string    	`json:"id" gorm:"primaryKey  binding:"required"`
	Username string 	`json:"username"  binding:"required"`
	Version  int  		`json:"version" gorm:"default:1"` // optimistic locking
}