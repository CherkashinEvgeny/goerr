package errors

import (
	"encoding/json"
	"encoding/xml"
)

type Code string

type Error struct {
	code      Code
	message   string
	paramsMap map[string]any
	Stack
}

func New(template Template, params ...Param) error {
	return newError(template, params)
}

func newError(template Template, params Params) Error {
	paramsMap := mergeParamMaps(template.Params.toMap(), params.toMap())
	var stack Stack
	if cfg.CollectStackTrace {
		stack = callers(4)
	}
	checkTemplateParams(template, paramsMap)
	message := format(template.Message, paramsMap)
	return Error{
		code:      template.Code,
		message:   message,
		paramsMap: paramsMap,
		Stack:     stack,
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

func checkTemplateParams(template Template, actualParams map[string]any) {
	templateParams := params(template.Message)
	for _, param := range templateParams {
		_, found := actualParams[param]
		if !found {
			cfg.MissingTemplateParamHandler(template, param)
		}
	}
}

func (e Error) Code() Code {
	return e.code
}

func (e Error) Error() string {
	return format(e.message, e.paramsMap)
}

func (e Error) Value(key string) any {
	if e.paramsMap == nil {
		return nil
	}
	return e.paramsMap[key]
}

func (e Error) Params() Params {
	params := make(Params, 0, len(e.paramsMap))
	for key, value := range e.paramsMap {
		params = append(params, Param{key, value})
	}
	return params
}

func (e Error) MarshalJSON() ([]byte, error) {
	fieldsCount := 2
	if cfg.MarshalStackTrace {
		fieldsCount++
	}
	data := make(map[string]any, len(e.paramsMap)+fieldsCount)
	for key, value := range e.paramsMap {
		if cfg.IsPrivateParam(key) {
			continue
		}
		data[cfg.JsonKey(key)] = value
	}
	data[cfg.JsonKey("message")] = e.message
	data[cfg.JsonKey("code")] = e.code
	if cfg.MarshalStackTrace {
		data[cfg.JsonKey("stackTrace")] = e.StackTrace().String()
	}
	return json.Marshal(data)
}

func (e Error) MarshalXML(en *xml.Encoder, start xml.StartElement) (err error) {
	fieldsCount := 2
	if cfg.MarshalStackTrace {
		fieldsCount++
	}
	data := make(map[string]any, len(e.paramsMap)+fieldsCount)
	for key, value := range e.paramsMap {
		if cfg.IsPrivateParam(key) {
			continue
		}
		data[cfg.XMLKey(key)] = value
	}
	if cfg.MarshalStackTrace {
		data[cfg.XMLKey("StackTrace")] = e.StackTrace().String()
	}
	data[cfg.XMLKey("Message")] = e.message
	data[cfg.XMLKey("Code")] = e.code
	start.Name.Local = "Error"
	return en.EncodeElement(data, start)
}
