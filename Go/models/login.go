package models

import "gorm.io/gorm"

type Login struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserID;references:ID"`
	UserID uint `json:"user_id"`
}
