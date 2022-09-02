package dto

// RegisterForm - dto для регистрации УЗ
type RegisterForm struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

// ErrorResponse - dto для ответа с ошибкой
type ErrorResponse struct {
	Status      bool   `json:"status"`
	Description string `json:"description"`
}
