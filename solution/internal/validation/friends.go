package validation

import "solution/internal/contract"

func ValidateLoginBody(body contract.LoginBody) error {
	if body.Login == "" {
		return contract.MISSING_FIELD("login")
	}
	return nil
}
