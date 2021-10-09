package models

type Response struct {
	Success bool   `json:"success"`
	Payload string `json:"payload"`
	Error   string `json:"error"`
}
