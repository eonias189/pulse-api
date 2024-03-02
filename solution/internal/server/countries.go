package server

import (
	"slices"
	"solution/internal/contract"
	"solution/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func validateRegion(region string, s *service.Service) error {
	if region == "" {
		return nil
	}
	regions, err := s.GetRegions()
	if err != nil {
		return err
	}
	if !slices.Contains(regions, region) {
		return contract.UNKNOWN_REGION(region)
	}
	return nil
}

func validateAlpha2(alpha2 string, s *service.Service) error {
	if len(alpha2) != 2 {
		return contract.BAD_PATH_PARAM(alpha2)
	}
	return nil
}

func handleCountries(r fiber.Router, s *service.Service) {
	r.Get("/", func(c *fiber.Ctx) error {
		region := c.Query("region")

		err := validateRegion(region, s)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}

		var (
			countries []contract.Country
		)
		if region == "" {
			countries, err = s.GetCountries()
		} else {
			countries, err = s.GetCountriesOfRegion(region)
		}
		if err != nil {
			return sendError(c, err, fiber.StatusInternalServerError)
		}
		slices.SortFunc(countries, func(c1 contract.Country, c2 contract.Country) int {
			return strings.Compare(c1.Alpha2, c2.Alpha2)
		})
		return c.JSON(countries)
	})

	r.Get("/:alpha2", func(c *fiber.Ctx) error {
		alpha2 := c.Params("alpha2")
		err := validateAlpha2(alpha2, s)
		if err != nil {
			return sendError(c, err, fiber.StatusBadRequest)
		}
		country, err := s.GetCountryByAlpha2(alpha2)
		if err != nil {
			return sendError(c, err, fiber.StatusNotFound)
		}
		return c.JSON(country)
	})
}
