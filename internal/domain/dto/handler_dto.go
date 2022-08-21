package dto

// RegisterForm - dto для регистрации УЗ
type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
