package errors

type Template struct {
	Code    Code
	Message string
	Params  Params
}

var NotFound = Template{
	Code:    "NOT_FOUND",
	Message: "{Resource} not found",
	Params: Params{
		Resource("Resource"),
	},
}

var AlreadyExists = Template{
	Code:    "ALREADY_EXISTS",
	Message: "{Resource} already exists",
	Params: Params{
		Resource("Resource"),
	},
}

var Unauthorized = Template{
	Code:    "UNAUTHORIZED",
	Message: "Unauthorized",
}

var Forbidden = Template{
	Code:    "FORBIDDEN",
	Message: "Forbidden",
}

var NotAllowed = Template{
	Code:    "NOT_ALLOWED",
	Message: "Not allowed: {Cause}",
	Params: Params{
		Cause("illegal state"),
	},
}

var ToManyRequests = Template{
	Code:    "TO_MANY_REQUESTS",
	Message: "To many request, please reduce your requests rate",
}

var Timeout = Template{
	Code:    "TIMEOUT",
	Message: "Timeout",
}

var Canceled = Template{
	Code:    "CANCELED",
	Message: "Operation canceled",
}

var InternalError = Template{
	Code:    "INTERNAL_ERROR",
	Message: "Internal error",
}

var NotImplemented = Template{
	Code:    "NOT_IMPLEMENTED",
	Message: "Not implemented",
}
