package dto

//dto для регистрации УЗ
type registerForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
