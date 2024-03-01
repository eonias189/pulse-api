package contract

type Country struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Region string `json:"region"`
}

type ErrorResp struct {
	Reason string `json:"reason"`
}

func NewErrorResp(err error) ErrorResp {
	return ErrorResp{Reason: err.Error()}
}
