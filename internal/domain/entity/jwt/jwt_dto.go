package jwt

type Header struct {
	Alg  string `json:"alg"`
	Type string `json:"type"`
}

type Payload struct {
	UserId string `json:"userId"`
	Exp    string `json:"exp"`
}
