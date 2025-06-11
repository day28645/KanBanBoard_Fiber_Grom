package models

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	User      User   `gorm:"foreignKey:OwnerID;references:ID"`
	BoardName string `gorm:"column:board_name;size:20;unique" json:"board_name"`
	OwnerID   uint   `json:"owner_id"`
}
