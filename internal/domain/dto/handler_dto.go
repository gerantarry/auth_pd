package dto

// RegisterRequestDto - dto для регистрации УЗ
type RegisterRequestDto struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

// StatusResponse - dto для ответа содержащего статус по проделанной операции
type StatusResponse struct {
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

// LoginRequestDto - dto запроса для операции логина клиента
type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
