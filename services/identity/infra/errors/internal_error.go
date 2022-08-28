package errors

import gerrors "errors"

type InternalError struct {
	err error
}

func NewInternalError(text string) InternalError {
	return InternalError{
		err: gerrors.New(text),
	}
}

func (r InternalError) Error() string {
	return r.err.Error()
}
