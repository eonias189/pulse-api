package contract

import (
	"fmt"
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
