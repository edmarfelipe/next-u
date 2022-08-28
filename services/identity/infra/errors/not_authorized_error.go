package errors

import gerrors "errors"

type NotAuthorizedError struct {
	err error
}

func NewNotAuthorizedError(text string) NotAuthorizedError {
	return NotAuthorizedError{
		err: gerrors.New(text),
	}
}

func (r NotAuthorizedError) Error() string {
	return r.err.Error()
}
