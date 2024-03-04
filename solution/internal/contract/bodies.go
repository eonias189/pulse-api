package contract

type RegisterBody User

type SignInBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MeProfileBody struct {
	CountryCode string `json:"countryCode"`
	IsPublic    bool   `json:"isPublic"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
}

func (mp MeProfileBody) ToUser(last User) User {
	return User{
		Login:           last.Login,
		Email:           last.Email,
		Password:        last.Password,
		CountryCode:     mp.CountryCode,
		IsPublic:        mp.IsPublic,
		Phone:           mp.Phone,
		Image:           mp.Image,
		PasswordChanged: last.PasswordChanged,
	}
}

type MeUpdatePasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type LoginBody struct {
	Login string `json:"login"`
}

type PostNewBody struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
