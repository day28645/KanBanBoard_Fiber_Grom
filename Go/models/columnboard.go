package models

import "gorm.io/gorm"

type ColumnBoard struct {
	gorm.Model
	Board      Board  `gorm:"foreignKey:BoardID;references:ID"`
	ColumnName string `gorm:"column:column_name;size:20;" json:"column_name"`
	BoardID    uint   `json:"board_id"`
}
