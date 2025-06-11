package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ColumnBoard    ColumnBoard `gorm:"foreignKey:ColumnBoardID;references:ID"`
	User           User        `gorm:"foreignKey:CreateByUserID;references:ID"`
	Title          string      `gorm:"column:title;size:20;" json:"title"`
	DueDate        time.Time   `json:"due_date"`
	ColumnBoardID  uint        `json:"column_board_id"`
	CreateByUserID uint        `json:"create_by_user_id"`
}
