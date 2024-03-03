package contract

import (
	"fmt"
)

var (
	INVALID_PASSWORD    = fmt.Errorf("invalid password")
	USER_ALREADY_EXISTS = fmt.Errorf("user with some of this data already exists")
	BAD_CRENDIALS       = fmt.Errorf("bad crendials")
	PASSWORD_CHANGED    = fmt.Errorf("password changed")
	ACCESS_DENIED       = fmt.Errorf("access denied")
	INCORRECT_PASSWORD  = fmt.Errorf("password is incorrect")
)

func ENV_ERROR(param string) error {
	return fmt.Errorf("missing %v in env", param)
}

func UNKNOWN_REGION(region string) error {
	return fmt.Errorf("unknown region: %v", region)
}

func BAD_PATH_PARAM(param string) error {
	return fmt.Errorf("bad path param: %v", param)
}

func MISSING_FIELD(field string) error {
	return fmt.Errorf("missing required field %v", field)
}

func BAD_BODY_PARAM(msg string) error {
	return fmt.Errorf("bad body param: %v", msg)
}

func UNKNOWN_COUNTRY_CODE(code string) error {
	return fmt.Errorf("unknown country code: %v", code)
}

func NOT_FOUND(what, which string) error {
	return fmt.Errorf("%v not found: %v", what, which)
}

func ALREADY_IN_FRIENDS(user string) error {
	return fmt.Errorf("user %v is already in friends", user)
}
