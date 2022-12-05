package errors

import gerrors "errors"

type InvalidInputError struct {
	err error
}

func NewInvalidInputError(text string) InvalidInputError {
	return InvalidInputError{
		err: gerrors.New(text),
	}
}

func (r InvalidInputError) Error() string {
	return r.err.Error()
}
