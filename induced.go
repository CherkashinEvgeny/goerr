package errors

import "fmt"

type InducedError struct {
	err   Error
	cause error
}

func newInducedError(err error, template Template, params Params) InducedError {
	return InducedError{
		err:   newError(template, params),
		cause: err,
	}
}

func (e InducedError) Code() Code {
	return e.err.Code()
}

func (e InducedError) Message() string {
	return e.err.Message()
}

func (e InducedError) Error() string {
	return fmt.Sprintf("%s: %s", e.err.Error(), e.cause.Error())
}

func (e InducedError) Cause() error {
	return e.cause
}

func (e InducedError) Trace() []string {
	return e.err.Trace()
}
