package contract

type RegisterBody User

type SignInBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
