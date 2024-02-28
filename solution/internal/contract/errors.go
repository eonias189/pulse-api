package contract

import (
	"fmt"
)

func ENV_ERROR(param string) error {
	return fmt.Errorf("missing %v in env", param)
}
