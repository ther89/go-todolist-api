package models

type TodoModel struct {
	Id          int    `gorm:"primary_key" json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
