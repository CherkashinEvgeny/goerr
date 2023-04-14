package errors

import "fmt"

const CodeValidationError Code = "ValidationError"

func ValidationError(params ...Param) error {
	return New(validationError, params...)
}

var validationError = Template{
	Code: CodeValidationError,
	Message: func(params map[string]any) string {
		resource, found := params[keyResource]
		if !found {
			return "Resource validation error"
		}
		return fmt.Sprintf("%s validation error", resource)
	},
	Params: Params{},
}

const CodeBlockingLink Code = "BlockingLink"

func BlockingLink(params ...Param) error {
	return New(blockingLink, params...)
}

var blockingLink = Template{
	Code: CodeBlockingLink,
	Message: func(params map[string]any) string {
		return "There are links to this resource"
	},
	Params: Params{},
}

const CodeChecksumError = "ChecksumError"

func ChecksumError(params ...Param) error {
	return New(checksumError, params...)
}

var checksumError = Template{
	Code: CodeChecksumError,
	Message: func(params map[string]any) string {
		return "Checksum does not match"
	},
	Params: Params{},
}

const CodeUnauthorized Code = "Unauthorized"

func Unauthorized(params ...Param) error {
	return New(unauthorized, params...)
}

var unauthorized = Template{
	Code: CodeUnauthorized,
	Message: func(params map[string]any) string {
		return fmt.Sprintf("Unauthorized")
	},
	Params: Params{},
}

const CodeForbidden Code = "Forbidden"

func Forbidden(params ...Param) error {
	return New(forbidden, params...)
}

var forbidden = Template{
	Code: CodeForbidden,
	Message: func(params map[string]any) string {
		return fmt.Sprintf("Forbidden")
	},
	Params: Params{},
}

const CodeNotFound Code = "NotFound"

func NotFound(params ...Param) error {
	return New(notFound, params...)
}

var notFound = Template{
	Code: CodeNotFound,
	Message: func(params map[string]any) string {
		resource, found := params[keyResource]
		if !found {
			return "Resource not found"
		}
		return fmt.Sprintf("%s not found", resource)
	},
	Params: Params{},
}

const CodeTimeout Code = "Timeout"

func Timeout(params ...Param) error {
	return New(timeout, params...)
}

var timeout = Template{
	Code:    CodeTimeout,
	Message: Message(`Timeout`),
	Params:  Params{},
}

const CodeAlreadyExists Code = "AlreadyExists"

func AlreadyExists(params ...Param) error {
	return New(alreadyExists, params...)
}

var alreadyExists = Template{
	Code: CodeAlreadyExists,
	Message: func(params map[string]any) string {
		resource, found := params[keyResource]
		if !found {
			return "Resource already exists"
		}
		return fmt.Sprintf("%s already exists", resource)
	},
	Params: Params{},
}

const CodeAlreadyInProgress Code = "AlreadyInProgress"

func AlreadyInProgress(params ...Param) error {
	return New(alreadyInProgress, params...)
}

var alreadyInProgress = Template{
	Code: CodeAlreadyInProgress,
	Message: func(params map[string]any) string {
		return "Already in progress"
	},
	Params: Params{},
}

const CodeIllegalState Code = "IllegalState"

func IllegalState(params ...Param) error {
	return New(illegalState, params...)
}

var illegalState = Template{
	Code: CodeIllegalState,
	Message: func(params map[string]any) string {
		reason, found := params[keyReason]
		if !found {
			return "Illegal state"
		}
		return fmt.Sprintf("Illegal state: %s", reason)
	},
	Params: Params{},
}

const CodePreconditionFailed Code = "PreconditionFailed"

func PreconditionFailed(params ...Param) error {
	return New(preconditionFailed, params...)
}

var preconditionFailed = Template{
	Code: CodePreconditionFailed,
	Message: func(params map[string]any) string {
		precondition, found := params[keyPrecondition]
		if !found {
			return "Precondition failed"
		}
		return fmt.Sprintf("Precondition %s failed", precondition)
	},
	Params: Params{},
}

const CodePreconditionRequired Code = "PreconditionRequired"

func PreconditionRequired(params ...Param) error {
	return New(preconditionRequired, params...)
}

var preconditionRequired = Template{
	Code: CodePreconditionFailed,
	Message: func(params map[string]any) string {
		precondition, found := params[keyPrecondition]
		if !found {
			return "Precondition required"
		}
		return fmt.Sprintf("Precondition %s required", precondition)
	},
	Params: Params{},
}

const CodeToManyRequests Code = "ToManyRequests"

func ToManyRequests(params ...Param) error {
	return New(toManyRequests, params...)
}

var toManyRequests = Template{
	Code:    CodeToManyRequests,
	Message: Message("To many request, please reduce your requests rate"),
}

const CodeInternalError Code = "InternalError"

func InternalError(err error) error {
	return Wrap(err, internalError)
}

var internalError = Template{
	Code: CodeInternalError,
	Message: func(params map[string]any) string {
		return "Internal error"
	},
	Params: Params{},
}

const CodeNotImplemented Code = "NotImplemented"

func NotImplemented(params ...Param) error {
	return New(notImplemented, params...)
}

var notImplemented = Template{
	Code: CodeNotImplemented,
	Message: func(params map[string]any) string {
		return "Not implemented"
	},
	Params: Params{},
}
