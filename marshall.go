package errors

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

var _ json.Marshaler = (*Error)(nil)

func (e Error) MarshalJSON() ([]byte, error) {
	fieldsCount := 2 + len(e.paramsMap)
	if cfg.MarshalStackTrace {
		fieldsCount++
	}
	var err error
	data := make(map[string]json.RawMessage, fieldsCount)
	for key, value := range e.paramsMap {
		if cfg.IsPrivateParam(key) {
			continue
		}
		data[cfg.ParamNameToJsonKey(key)], err = marshalJson(key, value)
		if err != nil {
			return nil, err
		}
	}
	data[cfg.ParamNameToJsonKey(keyCode)], err = marshalJson(keyCode, e.code)
	if err != nil {
		return nil, keyMarshalError{keyCode, err}
	}
	data[cfg.ParamNameToJsonKey(keyMessage)], err = marshalJson(keyMessage, e.message)
	if err != nil {
		return nil, keyMarshalError{keyMessage, err}
	}
	if cfg.MarshalCause && e.cause != nil {
		data[cfg.ParamNameToJsonKey(keyCause)], err = marshalJson(keyCause, e.cause)
		if err != nil {
			return nil, keyMarshalError{keyCause, err}
		}
	}
	if cfg.MarshalStackTrace && e.stackTrace != nil {
		data[cfg.ParamNameToJsonKey(keyStackTrace)], err = marshalJson(keyStackTrace, e.stackTrace)
		if err != nil {
			return nil, keyMarshalError{keyStackTrace, err}
		}
	}
	return json.Marshal(data)
}

func marshalJson(key string, value any) ([]byte, error) {
	marshaller, found := cfg.MarshalJsonParam[key]
	if !found {
		marshaller = cfg.MarshalJson
	}
	jsonValue, err := marshaller(value)
	if err != nil {
		return nil, err
	}
	return jsonValue, nil
}

var _ json.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalJSON(bytes []byte) error {
	data := map[string]json.RawMessage{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	codeJson, ok := data[cfg.ParamNameToJsonKey(keyCode)]
	delete(data, keyCode)
	if !ok {
		return codeIsMissingError
	}
	var code Code
	err = unmarshalJson(keyCode, codeJson, &code)
	if err != nil {
		return keyUnmarshalError{keyCode, err}
	}
	messageJson, ok := data[cfg.ParamNameToJsonKey(keyMessage)]
	delete(data, keyMessage)
	if !ok {
		return messageIsMissingError
	}
	var message string
	err = unmarshalJson(keyMessage, messageJson, message)
	if err != nil {
		return keyUnmarshalError{keyMessage, err}
	}
	causeJson, ok := data[cfg.ParamNameToJsonKey(keyCause)]
	delete(data, keyCause)
	var cause error
	if ok {
		err = unmarshalJson(keyCause, causeJson, &cause)
		if err != nil {
			return keyUnmarshalError{keyCause, err}
		}
	}
	delete(data, cfg.ParamNameToJsonKey(keyStackTrace))

	paramsMap := make(map[string]any, len(data))
	for jkey, jvalue := range data {
		key := cfg.ParamNameFromJsonKey(jkey)
		var value any
		err = unmarshalJson(key, jvalue, &value)
		if err != nil {
			return keyUnmarshalError{key, err}
		}
		paramsMap[key] = jvalue
	}

	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(1)
	}

	*e = Error{
		code:       code,
		message:    message,
		cause:      cause,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
	}
	return nil
}

func unmarshalJson(key string, data []byte, v any) error {
	unmarshller, found := cfg.UnmarshalJsonParam[key]
	if !found {
		unmarshller = cfg.UnmarshalJson
	}
	return unmarshller(data, v)
}

var _ xml.Marshaler = (*Error)(nil)

func (e Error) MarshalXML(en *xml.Encoder, start xml.StartElement) error {
	err := en.EncodeToken(start)
	if err != nil {
		return err
	}
	err = marshalXml(keyCode, en, e.code)
	if err != nil {
		return keyMarshalError{keyCode, err}
	}
	err = marshalXml(keyMessage, en, e.message)
	if err != nil {
		return keyMarshalError{keyMessage, err}
	}
	if cfg.MarshalCause && e.cause != nil {
		err = marshalXml(keyCause, en, e.cause)
		if err != nil {
			return keyMarshalError{keyCause, err}
		}
	}
	if cfg.MarshalStackTrace && e.stackTrace != nil {
		err = marshalXml(keyStackTrace, en, e.stackTrace)
		if err != nil {
			return keyMarshalError{keyStackTrace, err}
		}
	}

	for key, value := range e.paramsMap {
		err = marshalXml(key, en, value)
		if err != nil {
			return keyMarshalError{key, err}
		}
	}
	return en.EncodeToken(start.End())
}

func marshalXml(key string, en *xml.Encoder, v any) error {
	marshaller, found := cfg.MarshalXmlParam[key]
	if !found {
		marshaller = cfg.MarshalXml
	}
	return marshaller(en, xml.StartElement{Name: xml.Name{Local: cfg.ParamNameToXMLKey(key)}}, v)
}

var _ xml.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	var code Code
	var message string
	var cause error
	var paramsMap map[string]any
	for {
		token, _ := d.Token()
		if token == nil {
			break
		}
		el, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		key := cfg.ParamNameFromXMLKey(el.Name.Local)
		switch key {
		case keyCode:
			err := unmarshalXml(keyCode, d, el, &code)
			if err != nil {
				return keyUnmarshalError{keyCode, err}
			}
		case keyMessage:
			err := unmarshalXml(keyMessage, d, el, &message)
			if err != nil {
				return keyUnmarshalError{keyMessage, err}
			}
		case keyCause:
			err := unmarshalXml(keyCause, d, el, &cause)
			if err != nil {
				return keyUnmarshalError{keyCause, err}
			}
		case keyStackTrace:
			break
		default:
			var value any
			err := unmarshalXml(key, d, el, &value)
			if err != nil {
				return keyUnmarshalError{key, err}
			}
			paramsMap[key] = value
		}
	}

	var stackTrace StackTrace
	if cfg.CollectStackTrace {
		stackTrace = trace(1)
	}

	*e = Error{
		code:       code,
		message:    message,
		cause:      cause,
		paramsMap:  paramsMap,
		stackTrace: stackTrace,
	}
	return nil
}

func unmarshalXml(key string, d *xml.Decoder, start xml.StartElement, v any) error {
	unmarshaller, found := cfg.UnmarshalXmlParam[key]
	if !found {
		unmarshaller = cfg.UnmarshalXml
	}
	return unmarshaller(d, start, v)
}

var codeIsMissingError = errors.New("code is missing")

var messageIsMissingError = errors.New("message is missing")

type keyMarshalError struct {
	key string
	err error
}

func (e keyMarshalError) Error() string {
	return fmt.Sprintf("marshal %s: %s", e.key, e.err.Error())
}

func (e keyMarshalError) Unwrap() error {
	return e.err
}

func (e keyMarshalError) Cause() error {
	return e.err
}

type keyUnmarshalError struct {
	key string
	err error
}

func (e keyUnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal %s: %s", e.key, e.err.Error())
}

func (e keyUnmarshalError) Unwrap() error {
	return e.err
}

func (e keyUnmarshalError) Cause() error {
	return e.err
}
