package validation

import (
	"slices"
	"solution/internal/contract"
	"solution/internal/service"
)

func ValidateRegion(regions []string, s *service.Service) error {
	if len(regions) == 0 {
		return nil
	}
	validRegions := []string{
		"Europe", "Africa", "Americas", "Oceania", "Asia",
	}
	for _, region := range regions {
		if !slices.Contains(validRegions, region) {
			return contract.UNKNOWN_REGION(region)
		}
	}
	return nil
}

func ValidateAlpha2(alpha2 string) error {
	if len(alpha2) != 2 {
		return contract.BAD_PATH_PARAM(alpha2)
	}

	return nil
}
