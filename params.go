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

func Resource(resource string) Param {
	return Param{"Resource", resource}
}

func Cause(cause string) Param {
	return Param{"Cause", cause}
}
