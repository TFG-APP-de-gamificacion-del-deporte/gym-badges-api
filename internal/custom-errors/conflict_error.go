package custom_errors

import (
	"fmt"
)

var (
	Conflict = ConflictError{CustomError{Message: "Conflict"}}
)

type ConflictError struct {
	CustomError
}

func (e ConflictError) Error() string {
	return e.Message
}

// BuildConflictError Builds a ConflictError with the supplied message, and the corresponding code.
func BuildConflictError(message string, args ...any) ConflictError {
	if len(args) == 0 {
		return ConflictError{CustomError{Message: message}}
	}
	return ConflictError{CustomError{Message: fmt.Sprintf(message, args...)}}
}
