package custom_errors

import (
	"fmt"
)

var (
	NotFound = NotFoundError{CustomError{Message: "Not Found"}}
)

type NotFoundError struct {
	CustomError
}

func (e NotFoundError) Error() string {
	return e.Message
}

// BuildNotFoundError Builds a NotFoundError with the supplied message, and the corresponding code.
func BuildNotFoundError(message string, args ...any) NotFoundError {
	if len(args) == 0 {
		return NotFoundError{CustomError{Message: message}}
	}
	return NotFoundError{CustomError{Message: fmt.Sprintf(message, args...)}}
}
