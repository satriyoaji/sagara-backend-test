package models

type User struct {
	Id uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" gorm:"unique" validate:"required,email,min=8,max=50"`
	Password []byte `json:"password" validate:"required,min=6,max=50"`
}