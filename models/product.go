package models

type Product struct {
	Id uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name" validate:"required,min=3,max=50"`
	Quantity int `json:"quantity" validate:"required,number"`
	Price int `json:"price" validate:"required,number"`
	Image string `json:"image"`
}