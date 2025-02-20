package expression

// Структура выражения
type Expression struct {
	id     int     `json:"id"`
	exp    string  `json:"exp"`
	status string  `json:status`
	result float64 `json:result`
}
