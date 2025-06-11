package models

import "gorm.io/gorm"

type TaskAssignee struct {
	gorm.Model
	Task             Task `gorm:"foreignKey:TaskID;references:ID"`
	Assignee         User `gorm:"foreignKey:UserID;references:ID"`
	AssignedBy       User `gorm:"foreignKey:AssignedByUserID;references:ID"`
	TaskID           uint `json:"task_id"`
	UserID           uint `json:"user_id"`
	AssignedByUserID uint `json:"assigned_by_user_id"`
}
