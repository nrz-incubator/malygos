package errors

import "fmt"

type ConflictError struct {
	what string
	kind string
}

func NewConflictError(kind string, what string) *ConflictError {
	return &ConflictError{what: what, kind: kind}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s %s already exists", e.kind, e.what)
}

func IsConflict(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

type NotFoundError struct {
	what string
	kind string
}

func NewNotFoundError(kind string, what string) *NotFoundError {
	return &NotFoundError{what: what, kind: kind}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.kind, e.what)
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
