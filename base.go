package errors

import "fmt"

const CodeValidationError Code = "ValidationError"

var ValidationError = Template{
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

var BlockingLink = Template{
	Code: CodeBlockingLink,
	Message: func(params map[string]any) string {
		return "There are links to this resource"
	},
	Params: Params{},
}

const CodeChecksumError = "ChecksumError"

var ChecksumError = Template{
	Code: CodeChecksumError,
	Message: func(params map[string]any) string {
		return "Checksum does not match"
	},
	Params: Params{},
}

const CodeUnauthorized Code = "Unauthorized"

var Unauthorized = Template{
	Code: CodeUnauthorized,
	Message: func(params map[string]any) string {
		return fmt.Sprintf("Unauthorized")
	},
	Params: Params{},
}

const CodeForbidden Code = "Forbidden"

var Forbidden = Template{
	Code: CodeForbidden,
	Message: func(params map[string]any) string {
		return fmt.Sprintf("Forbidden")
	},
	Params: Params{},
}

const CodeNotFound Code = "NotFound"

var NotFound = Template{
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

var Timeout = Template{
	Code:    CodeTimeout,
	Message: Message(`Timeout`),
	Params:  Params{},
}

const CodeAlreadyExists Code = "AlreadyExists"

var AlreadyExists = Template{
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

var AlreadyInProgress = Template{
	Code: CodeAlreadyInProgress,
	Message: func(params map[string]any) string {
		return "Already in progress"
	},
	Params: Params{},
}

const CodeIllegalState Code = "IllegalState"

var IllegalState = Template{
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

var PreconditionFailed = Template{
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

var PreconditionRequired = Template{
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

var ToManyRequests = Template{
	Code:    CodeToManyRequests,
	Message: Message("To many request, please reduce your requests rate"),
}

const CodeInternalError Code = "InternalError"

var InternalError = Template{
	Code: CodeInternalError,
	Message: func(params map[string]any) string {
		return "Internal error"
	},
	Params: Params{},
}

const CodeNotImplemented Code = "NotImplemented"

var NotImplemented = Template{
	Code: CodeNotImplemented,
	Message: func(params map[string]any) string {
		return "Not implemented"
	},
	Params: Params{},
}
