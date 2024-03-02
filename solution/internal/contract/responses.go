package contract

type ErrorResp struct {
	Reason string `json:"reason"`
}

func NewErrorResp(err error) ErrorResp {
	return ErrorResp{Reason: err.Error()}
}

type CountryResponse []Country

type RegisterResp struct {
	Profile UserProfile `json:"profile"`
}

type SignInResp struct {
	Token string `json:"token"`
}
