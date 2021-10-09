package models

type Todo struct {
	Id          int    `gorm:"primary_key" json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
