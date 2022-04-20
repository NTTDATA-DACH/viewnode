package srv

import (
	"fmt"
	"strings"
)

type UnauthorizedError struct {
	err error
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.err)
}

type NodesIsForbiddenError struct {
	err error
}

func (e NodesIsForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %v", e.err)
}

func DecorateError(err error) error {
	switch msg := err.Error(); {
	case msg == "Unauthorized":
		return UnauthorizedError{err: err}
	case strings.HasPrefix(msg, "nodes is forbidden"):
		return NodesIsForbiddenError{err: err}
	default:
		return err
	}
}