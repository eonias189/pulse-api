package validation

import "solution/internal/contract"

func ValidateNewPostBody(body contract.PostNewBody) error {
	if body.Content == "" {
		return contract.MISSING_FIELD("content")
	}

	if len(body.Tags) == 0 {
		return contract.MISSING_FIELD("tags")
	}

	return nil
}
