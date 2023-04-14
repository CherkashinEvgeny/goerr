package http

import (
	errors "github.com/CherkashinEvgeny/goerr"
	"net/http"
)

func Status(err error) int {
	e, ok := err.(errors.Error)
	if !ok {
		return http.StatusInternalServerError
	}
	status, found := GetStatus(e)
	if found {
		return status
	}
	switch e.Code() {
	case errors.CodeValidationError:
		return http.StatusBadRequest
	case errors.CodeBlockingLink:
		return http.StatusBadRequest
	case errors.CodeChecksumError:
		return http.StatusBadRequest
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeTimeout:
		return http.StatusRequestTimeout
	case errors.CodeAlreadyExists:
		return http.StatusConflict
	case errors.CodeAlreadyInProgress:
		return http.StatusConflict
	case errors.CodeIllegalState:
		return http.StatusConflict
	case errors.CodePreconditionFailed:
		return http.StatusPreconditionFailed
	case errors.CodePreconditionRequired:
		return http.StatusPreconditionRequired
	case errors.CodeToManyRequests:
		return http.StatusTooManyRequests
	case errors.CodeInternalError:
		return http.StatusInternalServerError
	case errors.CodeNotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}

const keyStatus = "httpStatus"

func WithStatus(status int) errors.Param {
	return errors.Param{keyStatus, status}
}

func GetStatus(e errors.Error) (int, bool) {
	val := e.Get(keyStatus)
	if val == nil {
		return 0, false
	}
	status, ok := val.(int)
	return status, ok
}
