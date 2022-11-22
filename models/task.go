package models

type Task struct {
	Id          string
	Name        string `json:"name"  gorm:"default:null" binding:"required"`
	Description string `json:"description,omitempty" binding:"max=3000"`
	Priority    string `json:"type" binding:"required"`
	UserID      uint   `gorm:"column:user_id"`
}
