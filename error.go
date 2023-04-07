package errors

func New(template Template, params Params) error {
	return newError(template, params)
}

func Wrap(err error, template Template, params Params) error {
	return newInducedError(err, template, params)
}
