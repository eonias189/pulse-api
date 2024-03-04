package validation

import (
	"solution/internal/contract"
	"strings"
)

func ValidateNewPostBody(bodyRaw string) error {
	if !strings.Contains(string(bodyRaw), `"content"`) {
		return contract.MISSING_FIELD("content")
	}

	if !strings.Contains(string(bodyRaw), `"tags"`) {
		return contract.MISSING_FIELD("tags")
	}

	return nil
}
