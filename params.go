package errors

type Params []Param

func (p Params) toMap() map[string]any {
	m := make(map[string]any, len(p))
	for _, param := range p {
		m[param.Name] = param.Value
	}
	return m
}

type Param struct {
	Name  string
	Value any
}

const keyId = "Id"

func WithId(id string) Param {
	return Param{keyId, id}
}

func Id(e Error) string {
	res, _ := e.Get(keyId).(string)
	return res
}

const keyResource = "Resource"

func WithResource(resource string) Param {
	return Param{keyResource, resource}
}

func Resource(e Error) string {
	res, _ := e.Get(keyResource).(string)
	return res
}

const keyCause = "Cause"

func WithCause(err error) Param {
	return Param{keyCause, err}
}

func Cause(e Error) error {
	err, _ := e.Get(keyCause).(error)
	return err
}
