package http

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	errors "github.com/CherkashinEvgeny/goerr"
)

func Status(err error) int {
	e, ok := err.(*errors.Error)
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

func GetStatus(err error) (int, bool) {
	e, ok := err.(*errors.Error)
	if !ok {
		return 0, false
	}
	status, ok := e.Get(keyStatus).(int)
	return status, ok
}

func init() {
	errors.Configure(func(config *errors.Config) {
		config.MarshalJsonParam[keyStatus] = func(value any) ([]byte, error) {
			return json.Marshal(value)
		}
		config.UnmarshalJsonParam[keyStatus] = func(data []byte) (any, error) {
			var status int
			err := json.Unmarshal(data, &status)
			if err != nil {
				return nil, err
			}
			return status, nil
		}
		config.MarshalXmlParam[keyStatus] = func(en *xml.Encoder, start xml.StartElement, value any) error {
			return en.EncodeElement(value, start)
		}
		config.UnmarshalXmlParam[keyStatus] = func(d *xml.Decoder, start xml.StartElement) (any, error) {
			var status int
			err := d.DecodeElement(&status, &start)
			if err != nil {
				return nil, err
			}
			return status, nil
		}
	})
}
