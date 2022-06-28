package user

type User struct {
	ID         int    `json:"ID"`
	FirstName  string `json:"FirstName"`
	SecondName string `json:"SecondName"`
	Login      string `json:"Login"`
	Password   string `json:"Password"`
	Email      string `json:"Email"`
}
