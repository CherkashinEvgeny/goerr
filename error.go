package errors

type Code string

const (
	keyCode       = "Code"
	keyMessage    = "Message"
	keyCause      = "Cause"
	keyStackTrace = "StackTrace"
)

var _ error = (*Error)(nil)

type Error struct {
	code       Code
	message    string
	cause      error
	paramsMap  map[string]any
	stackTrace StackTrace
}

func New(template Template, params ...Param) error {
	return newError(template, nil, params)
}

func Wrap(err error, template Template, params ...Param) error {
	return newError(template, err, params)
}

func Is(err error, template Template) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code() == template.Code
}

func newError(template Template, cause error, params Params) *Error {
	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(2)
	}
	paramsMap := mergeParamMaps(template.Params.toMap(), params.toMap())
	message := template.Message(paramsMap)
	code := template.Code
	return &Error{
		code:       code,
		message:    message,
		cause:      cause,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
	}
}

func mergeParamMaps(maps ...map[string]any) map[string]any {
	size := 0
	for _, p := range maps {
		size += len(p)
	}
	mergedParams := make(map[string]any, size)
	for _, m := range maps {
		for key, value := range m {
			mergedParams[key] = value
		}
	}
	return mergedParams
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Get(key string) any {
	switch key {
	case keyCode:
		return e.code
	case keyMessage:
		return e.message
	case keyCause:
		return e.cause
	case keyStackTrace:
		return e.StackTrace()
	default:
		if e.paramsMap == nil {
			return nil
		}
		return e.paramsMap[key]
	}
}

func (e *Error) Params() Params {
	params := make(Params, 0, len(e.paramsMap))
	for key, value := range e.paramsMap {
		params = append(params, Param{key, value})
	}
	return params
}

func (e *Error) StackTrace() StackTrace {
	return e.stackTrace
}
