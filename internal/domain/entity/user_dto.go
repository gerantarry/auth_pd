package entity

type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}
