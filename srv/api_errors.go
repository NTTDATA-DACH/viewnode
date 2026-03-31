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

func (e UnauthorizedError) Unwrap() error {
	return e.err
}

type NodesIsForbiddenError struct {
	err error
}

func (e NodesIsForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %v", e.err)
}

func (e NodesIsForbiddenError) Unwrap() error {
	return e.err
}

type ScopedEOFError struct {
	err error
}

func (e ScopedEOFError) Error() string {
	return fmt.Sprintf("scoped eof: %v", e.err)
}

func (e ScopedEOFError) Unwrap() error {
	return e.err
}

func DecorateError(err error) error {
	switch msg := err.Error(); {
	case msg == "Unauthorized":
		return UnauthorizedError{err: err}
	case strings.HasPrefix(msg, "nodes is forbidden"):
		return NodesIsForbiddenError{err: err}
	case strings.Contains(msg, "Get ") && strings.HasSuffix(msg, ": EOF"):
		return ScopedEOFError{err: err}
	default:
		return err
	}
}
