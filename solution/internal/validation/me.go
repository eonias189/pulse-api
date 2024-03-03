package validation

import "solution/internal/contract"

func ValidateUpdatePassword(body contract.MeUpdatePasswordBody) error {
	switch "" {
	case body.OldPassword:
		return contract.MISSING_FIELD("oldPassword")
	case body.NewPassword:
		return contract.MISSING_FIELD("newPassword")
	}
	return nil
}
