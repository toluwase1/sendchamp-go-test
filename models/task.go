package models

import "github.com/google/uuid"

type Task struct {
	Id          uuid.UUID
	Name        string `json:"name"  gorm:"unique;default:null" binding:"required"`
	Description string `json:"description,omitempty" binding:"max=3000"`
	Type        string `json:"type" binding:"required"`
}
