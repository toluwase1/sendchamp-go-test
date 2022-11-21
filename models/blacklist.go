package models

type BlackList struct {
	Model
	Token string `json:"token"`
	Email string `json:"email"`
}
