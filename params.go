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

func GetId(e Error) (string, bool) {
	val := e.Get(keyId)
	if val == nil {
		return "", false
	}
	id, ok := e.Get(keyId).(string)
	return id, ok
}

const keyResource = "Resource"

func WithResource(resource string) Param {
	return Param{keyResource, resource}
}

func GetResource(e Error) (string, bool) {
	val := e.Get(keyResource)
	if val == nil {
		return "", false
	}
	resource, ok := e.Get(keyResource).(string)
	return resource, ok
}

const keyReason = "Reason"

func WithReason(reason string) Param {
	return Param{keyReason, reason}
}

func GetReason(e Error) (string, bool) {
	val := e.Get(keyReason)
	if val == nil {
		return "", false
	}
	reason, ok := e.Get(keyReason).(string)
	return reason, ok
}

const keyValidationErrors = "Errors"

func WithValidationErrors(errors map[string]string) Param {
	return Param{keyValidationErrors, errors}
}

func GetValidationErrors(e Error) (map[string]string, bool) {
	val := e.Get(keyValidationErrors)
	if val == nil {
		return nil, false
	}
	errors, ok := e.Get(keyValidationErrors).(map[string]string)
	return errors, ok
}

const keyPrecondition = "Precondition"

func WithPrecondition(precondition string) Param {
	return Param{keyPrecondition, precondition}
}

func GetPrecondition(e Error) (string, bool) {
	val := e.Get(keyPrecondition)
	if val != nil {
		return "", false
	}
	precondition, ok := e.Get(keyPrecondition).(string)
	return precondition, ok
}
