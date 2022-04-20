package srv

import "fmt"

type UnauthorizedError struct {
	err error
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized error: %v", e.err)
}

