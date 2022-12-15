package httperror

import "errors"

type Error struct {
	Status int
	Err    error
}

func New(status int, message string) error {
	return &Error{
		Status: status,
		Err:    errors.New(message),
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
