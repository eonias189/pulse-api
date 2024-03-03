package server

import (
	"slices"
	"solution/internal/contract"
	"solution/internal/service"
	"solution/internal/utils"
	"solution/internal/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getRegions(c *fiber.Ctx) []string {
	regionsBytes := c.Request().URI().QueryArgs().PeekMulti("region")
	regions := utils.Map(regionsBytes, func(b []byte) string { return string(b) })
	regions = utils.Filter(regions, func(region string) bool { return region != "" })
	return regions
}

func handleCountries(r fiber.Router, s *service.Service) {
	r.Get("/", func(c *fiber.Ctx) error {
		regions := getRegions(c)

		err := validation.ValidateRegion(regions, s)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}

		var (
			countries []contract.Country
		)
		if len(regions) == 0 {
			countries, err = s.GetCountries()
		} else {
			countries, err = s.GetCountriesOfRegions(regions)
		}
		if err != nil {
			return utils.SendError(c, err, fiber.StatusInternalServerError)
		}
		slices.SortFunc(countries, func(c1 contract.Country, c2 contract.Country) int {
			return strings.Compare(c1.Alpha2, c2.Alpha2)
		})
		return c.JSON(countries)
	})

	r.Get("/:alpha2", func(c *fiber.Ctx) error {
		alpha2 := c.Params("alpha2")

		err := validation.ValidateAlpha2(alpha2)
		if err != nil {
			return utils.SendError(c, err, fiber.StatusBadRequest)
		}
		country, err := s.GetCountryByAlpha2(alpha2)
		if err != nil {
			return utils.SendError(c, contract.NOT_FOUND("country", alpha2), fiber.StatusNotFound)
		}
		return c.JSON(country)
	})
}
