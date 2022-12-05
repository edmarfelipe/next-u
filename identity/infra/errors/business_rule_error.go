package errors

import gerrors "errors"

type BusinessRuleError struct {
	err error
}

func NewBusinessRuleError(err string) BusinessRuleError {
	return BusinessRuleError{
		err: gerrors.New(err),
	}
}

func (r BusinessRuleError) Error() string {
	return r.err.Error()
}
