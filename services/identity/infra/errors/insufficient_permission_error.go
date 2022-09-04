package errors

import gerrors "errors"

type InsufficientPermissionError struct {
	err error
}

func NewInsufficientPermissionError(err string) InsufficientPermissionError {
	return InsufficientPermissionError{
		err: gerrors.New(err),
	}
}

func (r InsufficientPermissionError) Error() string {
	return r.err.Error()
}
