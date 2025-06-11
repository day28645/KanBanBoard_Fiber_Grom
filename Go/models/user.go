package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	IDCard   string `gorm:"column:id_card;size:13;unique" json:"id_card"`
	Username string `gorm:"column:username;size:50;unique" json:"username"`
	Password string `gorm:"column:password;size:10" json:"password"`
	Email    string `gorm:"column:email;size:20" json:"email"`
}
