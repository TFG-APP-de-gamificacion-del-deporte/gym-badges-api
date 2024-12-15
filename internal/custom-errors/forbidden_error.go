package custom_errors

import (
	"fmt"
)

var (
	Forbidden = ForbiddenError{CustomError{Message: "Forbidden"}}
)

type ForbiddenError struct {
	CustomError
}

func (e ForbiddenError) Error() string {
	return e.Message
}

// BuildNotFoundError Builds a NotFoundError with the supplied message, and the corresponding code.
func BuildForbiddenError(message string, args ...any) ForbiddenError {
	if len(args) == 0 {
		return ForbiddenError{CustomError{Message: message}}
	}
	return ForbiddenError{CustomError{Message: fmt.Sprintf(message, args...)}}
}
