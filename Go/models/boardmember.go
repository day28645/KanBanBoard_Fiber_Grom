package models

import (
	"gorm.io/gorm"
)

type BoardMember struct {
	gorm.Model
	User    User
	Board   Board
	Role    string `gorm:"column:role;size:20;" json:"role"`
	BoardID uint   `json:"board_id"`
	UserID  uint   `json:"user_id"`
}
