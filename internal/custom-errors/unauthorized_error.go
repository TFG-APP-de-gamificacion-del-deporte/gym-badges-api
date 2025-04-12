package custom_errors

import (
	"fmt"
)

var (
	Unauthorized = UnauthorizedError{CustomError{Message: "Unauthorized"}}
)

type UnauthorizedError struct {
	CustomError
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

// BuildUnauthorizedError Builds a UnauthorizedError with the supplied message, and the corresponding code.
func BuildUnauthorizedError(message string, args ...any) UnauthorizedError {
	if len(args) == 0 {
		return UnauthorizedError{CustomError{Message: message}}
	}
	return UnauthorizedError{CustomError{Message: fmt.Sprintf(message, args...)}}
}
