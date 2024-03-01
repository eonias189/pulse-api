package server

import (
	"slices"
	"solution/internal/contract"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (s *Server) handlePing(c fiber.Ctx) error {
	return c.SendString("ok")
}

func (s *Server) handleCountriesIndex(c fiber.Ctx) error {
	region := c.Query("region")

	err := s.validateRegion(region)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(contract.NewErrorResp(err))
	}

	var (
		countries []contract.Country
	)
	if region == "" {
		countries, err = s.service.GetCountries()
	} else {
		countries, err = s.service.GetCountriesOfRegion(region)
	}
	if err != nil {
		return c.JSON(contract.NewErrorResp(err))
	}
	slices.SortFunc(countries, func(c1 contract.Country, c2 contract.Country) int {
		return strings.Compare(c1.Alpha2, c2.Alpha2)
	})
	return c.JSON(countries)
}

func (s *Server) validateRegion(region string) error {
	if region == "" {
		return nil
	}
	regions, err := s.service.GetRegions()
	if err != nil {
		return err
	}
	if !slices.Contains(regions, region) {
		return contract.UNKNOWN_REGION(region)
	}
	return nil
}

func (s *Server) handleCountriesAlpha2(c fiber.Ctx) error {
	alpha2 := c.Params("alpha2")
	err := validateAlpha2(alpha2)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(contract.NewErrorResp(err))
	}
	country, err := s.service.GetCountryByAlpha2(alpha2)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(contract.NewErrorResp(err))
	}
	return c.JSON(country)
}

func validateAlpha2(alpha2 string) error {
	if len(alpha2) != 2 {
		return contract.BAD_PATH_PARAM(alpha2)
	}
	return nil
}
