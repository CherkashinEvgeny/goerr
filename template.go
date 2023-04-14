package errors

import (
	"strings"
	"text/template"
)

type Template struct {
	Code    Code
	Message func(params map[string]any) string
	Params  Params
}

func Message(str string) func(params map[string]any) string {
	tmpt, err := template.New("").Funcs(template.FuncMap{
		"default": func(defaultVal any, val any) any {
			if val == nil {
				return defaultVal
			}
			return val
		},
	}).Parse(str)
	if err != nil {
		panic(err)
	}
	return func(params map[string]any) string {
		sb := strings.Builder{}
		err = tmpt.Execute(&sb, params)
		if err != nil {
			return str
		}
		return sb.String()
	}
}

const CodeNotFound Code = "NotFound"

var NotFound = Template{
	Code:    CodeNotFound,
	Message: Message(`{{ .Resource | default "Resource" }} not found`),
	Params:  Params{},
}

const CodeAlreadyExists Code = "AlreadyExists"

var AlreadyExists = Template{
	Code:    CodeAlreadyExists,
	Message: Message(`{{ .Resource | default "Resource" }} already exists`),
	Params:  Params{},
}

const CodeUnauthorized Code = "Unauthorized"

var Unauthorized = Template{
	Code:    CodeUnauthorized,
	Message: Message(`Unauthorized`),
	Params:  Params{},
}

const CodeForbidden Code = "Forbidden"

var Forbidden = Template{
	Code:    CodeForbidden,
	Message: Message(`Forbidden`),
	Params:  Params{},
}

const CodeNotAllowed Code = "NotAllowed"

var NotAllowed = Template{
	Code:    CodeNotAllowed,
	Message: Message(`Not allowed{{ if ne .Cause nil }}: {{ .Cause }}{{ end }}`),
	Params:  Params{},
}

const CodeToManyRequests Code = "ToManyRequests"

var ToManyRequests = Template{
	Code:    CodeToManyRequests,
	Message: Message("To many request, please reduce your requests rate"),
}

const CodeTimeout Code = "Timeout"

var Timeout = Template{
	Code:    CodeTimeout,
	Message: Message(`Timeout`),
	Params:  Params{},
}

const CodeCanceled Code = "Canceled"

var Canceled = Template{
	Code:    CodeCanceled,
	Message: Message(`Operation canceled`),
	Params:  Params{},
}

const CodeInternalError Code = "InternalError"

func InternalError(err error) error {
	return Wrap(err, internalError)
}

var internalError = Template{
	Code:    CodeInternalError,
	Message: Message(`Internal error{{ if ne .Cause nil }}: {{ .Cause }}{{ end }}`),
	Params:  Params{},
}

const CodeNotImplemented Code = "NotImplemented"

var NotImplemented = Template{
	Code:    CodeNotImplemented,
	Message: Message(`Not implemented`),
	Params:  Params{},
}
