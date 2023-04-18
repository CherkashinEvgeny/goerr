package errors

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Code string

const (
	keyCode       = "Code"
	keyMessage    = "Message"
	keyStackTrace = "StackTrace"
)

var _ error = (*Error)(nil)

type Error struct {
	code       Code
	message    string
	stackTrace StackTrace
	paramsMap  map[string]any
}

func New(template Template, params ...Param) error {
	return newError(template, params)
}

func Wrap(err error, template Template, params ...Param) error {
	return newError(template, append(params, WithCause(err)))
}

func Is(err error) bool {
	_, ok := err.(Error)
	return ok
}

func Cast(err error) Error {
	e, _ := err.(Error)
	return e
}

func newError(template Template, params Params) Error {
	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(2)
	}
	paramsMap := mergeParamMaps(template.Params.toMap(), params.toMap())
	message := template.Message(paramsMap)
	code := template.Code
	return Error{
		code:       code,
		message:    message,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
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

func (e Error) Code() Code {
	return e.code
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Get(key string) any {
	switch key {
	case keyCode:
		return e.code
	case keyMessage:
		return e.message
	case keyStackTrace:
		return e.StackTrace()
	default:
		if e.paramsMap == nil {
			return nil
		}
		return e.paramsMap[key]
	}
}

func (e Error) Params() Params {
	params := make(Params, 0, len(e.paramsMap))
	for key, value := range e.paramsMap {
		params = append(params, Param{key, value})
	}
	return params
}

func (e Error) StackTrace() StackTrace {
	return e.stackTrace
}

var _ json.Marshaler = (*Error)(nil)

func (e Error) MarshalJSON() ([]byte, error) {
	fieldsCount := 2 + len(e.paramsMap)
	if cfg.MarshalStackTrace {
		fieldsCount++
	}
	data := make(map[string]any, fieldsCount)
	for key, value := range e.paramsMap {
		if cfg.IsPrivateParam(key) {
			continue
		}
		data[cfg.ToJsonKey(key)] = value
	}
	data[cfg.ToJsonKey(keyCode)] = e.code
	data[cfg.ToJsonKey(keyMessage)] = e.message
	if cfg.MarshalStackTrace && e.stackTrace != nil {
		data[cfg.ToJsonKey(keyStackTrace)] = stackTraceToStringArray(e.stackTrace)
	}
	return json.Marshal(data)
}

func stackTraceToStringArray(stackTrace StackTrace) []string {
	strs := make([]string, 0, len(stackTrace))
	for _, frame := range stackTrace {
		strs = append(strs, fmt.Sprintf("%s %s:%d", frame.Func(), frame.File(), frame.Line()))
	}
	return strs
}

var _ json.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalJSON(bytes []byte) error {
	data := map[string]any{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	code, ok := data[cfg.ToJsonKey(keyCode)].(string)
	delete(data, keyCode)
	if !ok {
		return codeIsMissingError
	}
	message, ok := data[cfg.ToJsonKey(keyMessage)].(string)
	delete(data, keyMessage)
	if !ok {
		return messageIsMissingError
	}
	delete(data, cfg.ToJsonKey(keyStackTrace))
	paramsMap := make(map[string]any, len(data))
	for key, value := range paramsMap {
		paramsMap[cfg.FromJsonKey(key)] = value
	}
	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(1)
	}
	*e = Error{
		code:       Code(code),
		message:    message,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
	}
	return nil
}

var _ xml.Marshaler = (*Error)(nil)

func (e Error) MarshalXML(en *xml.Encoder, start xml.StartElement) error {
	fieldsCount := 2 + len(e.paramsMap)
	if cfg.MarshalStackTrace {
		fieldsCount++
	}
	data := make(map[string]any, fieldsCount)
	for key, value := range e.paramsMap {
		if cfg.IsPrivateParam(key) {
			continue
		}
		data[cfg.ToXMLKey(key)] = value
	}
	data[cfg.ToXMLKey(keyCode)] = e.code
	data[cfg.ToXMLKey(keyMessage)] = e.message
	if cfg.MarshalStackTrace && e.stackTrace != nil {
		data[cfg.ToXMLKey(keyStackTrace)] = stackTraceToString(e.stackTrace)
	}
	start.Name.Local = "Error"
	return en.EncodeElement(data, start)
}

func stackTraceToString(stackTrace StackTrace) string {
	sb := strings.Builder{}
	for index, frame := range stackTrace {
		if index != 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(frame.Func())
		sb.WriteString("\n\t")
		sb.WriteString(frame.File())
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(frame.Line()))
	}
	return sb.String()
}

var _ xml.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	data := map[string]any{}
	err := d.Decode(&data)
	if err != nil {
		return err
	}
	code, ok := data[cfg.ToXMLKey(keyCode)].(string)
	delete(data, keyCode)
	if !ok {
		return codeIsMissingError
	}
	message, ok := data[cfg.ToXMLKey(keyMessage)].(string)
	delete(data, keyMessage)
	if !ok {
		return messageIsMissingError
	}
	delete(data, cfg.ToXMLKey(keyStackTrace))
	paramsMap := make(map[string]any, len(data))
	for key, value := range paramsMap {
		paramsMap[cfg.FromXMLKey(key)] = value
	}
	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(1)
	}
	*e = Error{
		code:       Code(code),
		message:    message,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
	}
	return nil
}

var codeIsMissingError = errors.New("code is missing")

var messageIsMissingError = errors.New("message is missing")
