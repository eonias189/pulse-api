package validation

import (
	"solution/internal/contract"
	"solution/internal/service"
	"strings"
)

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return contract.INVALID_PASSWORD
	}

	latin := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"

	if !strings.ContainsAny(password, latin) {
		return contract.INVALID_PASSWORD
	}

	if !strings.ContainsAny(password, strings.ToUpper(latin)) {
		return contract.INVALID_PASSWORD
	}

	if !strings.ContainsAny(password, digits) {
		return contract.INVALID_PASSWORD
	}

	return nil

}

func ValidateImage(image string) error {
	if len(image) > 200 {
		return contract.BAD_BODY_PARAM("image len must be not longer than 200 symbols")
	}
	return nil
}

func ValidateRegister(body contract.RegisterBody, s *service.Service) error {
	switch "" {
	case body.Login:
		return contract.MISSING_FIELD("login")
	case body.Email:
		return contract.MISSING_FIELD("email")
	case body.Password:
		return contract.MISSING_FIELD("password")
	case body.CountryCode:
		return contract.MISSING_FIELD("countryCode")
	}

	err := ValidateImage(body.Image)
	if err != nil {
		return err
	}

	err = ValidateAlpha2(body.CountryCode)
	if err != nil {
		return err
	}

	_, err = s.GetCountryByAlpha2(body.CountryCode)
	if err != nil {
		return contract.UNKNOWN_COUNTRY_CODE(body.CountryCode)
	}

	err = ValidatePassword(body.Password)
	if err != nil {
		return err
	}

	return nil
}

func ValidateSignIn(body contract.SignInBody, s *service.Service) error {
	switch "" {
	case body.Login:
		return contract.MISSING_FIELD("login")
	case body.Password:
		return contract.MISSING_FIELD("password")
	}
	return nil
}
