package models

// Структура выражения
type Expression struct {
	Id     int     `json:"id"`
	Exp    string  `json:"exp"`
	Status string  `json:status`
	Result float64 `json:result`
}
