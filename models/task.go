package models

type Task struct {
	Id          string
	Name        string `json:"name"  gorm:"unique;default:null" binding:"required"`
	Description string `json:"description,omitempty" binding:"max=3000"`
	Priority    string `json:"type" binding:"required"`
}
