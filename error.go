package errors

type Code string

type Error struct {
	template Template
	Params   map[string]any
	err      error
	trace    []string
}

func New(template Template, params map[string]any) error {
	return Error{
		template: template,
		Params:   params,
		err:      nil,
		trace:    Traces(3),
	}
}

func Wrap(err error, template Template, params map[string]any) error {
	return Error{
		template: template,
		Params:   params,
		err:      err,
		trace:    Traces(3),
	}
}

func (e Error) Code() Code {
	return e.template.Code
}

func (e Error) Message() string {
	return e.template.Format(e.Params)
}

func (e Error) Cause() error {
	return e.err
}

func (e Error) Error() string {
	if e.err == nil {
		return e.Message()
	}
	return e.err.Error()
}

func (e Error) Trace() []string {
	return e.trace
}
