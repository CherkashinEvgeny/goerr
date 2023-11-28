package errors

import (
	"encoding/json"
	"encoding/xml"
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
		data[cfg.MarshalJsonKey(key)], err = marshalJson(key, value)
		if err != nil {
			return nil, err
		}
	}
	data[cfg.MarshalJsonKey(keyCode)], err = marshalJson(keyCode, e.code)
	if err != nil {
		return nil, keyMarshalError{keyCode, err}
	}
	data[cfg.MarshalJsonKey(keyMessage)], err = marshalJson(keyMessage, e.message)
	if err != nil {
		return nil, keyMarshalError{keyMessage, err}
	}
	if cfg.MarshalCause && e.cause != nil {
		data[cfg.MarshalJsonKey(keyCause)], err = marshalJson(keyCause, e.cause)
		if err != nil {
			return nil, keyMarshalError{keyCause, err}
		}
	}
	if cfg.MarshalStackTrace && e.stackTrace != nil {
		data[cfg.MarshalJsonKey(keyStackTrace)], err = marshalJson(keyStackTrace, e.stackTrace)
		if err != nil {
			return nil, keyMarshalError{keyStackTrace, err}
		}
	}
	return json.Marshal(data)
}

func marshalJson(key string, value any) ([]byte, error) {
	marshaller, found := cfg.MarshalJsonParam[key]
	if found {
		return marshaller(value)
	}
	jsonValue, err := json.Marshal(value)
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
	codeJson, ok := data[cfg.MarshalJsonKey(keyCode)]
	delete(data, keyCode)
	if !ok {
		return keyMissingError{keyCode}
	}
	codeValue, err := unmarshalJson(keyCode, codeJson)
	if err != nil {
		return keyUnmarshalError{keyCode, err}
	}
	code, ok := codeValue.(Code)
	if !ok {
		return keyCastError{keyCode}
	}
	messageJson, ok := data[cfg.MarshalJsonKey(keyMessage)]
	delete(data, keyMessage)
	if !ok {
		return keyMissingError{keyMessage}
	}
	messageValue, err := unmarshalJson(keyMessage, messageJson)
	if err != nil {
		return keyUnmarshalError{keyMessage, err}
	}
	message, ok := messageValue.(string)
	if !ok {
		return keyCastError{keyMessage}
	}
	causeJson, ok := data[cfg.MarshalJsonKey(keyCause)]
	var cause error
	if ok {
		delete(data, keyCause)
		causeValue, err := unmarshalJson(keyCause, causeJson)
		if err != nil {
			return keyUnmarshalError{keyCause, err}
		}
		cause, ok = causeValue.(error)
		if !ok {
			return keyCastError{keyCause}
		}
	}
	delete(data, cfg.MarshalJsonKey(keyStackTrace))

	paramsMap := make(map[string]any, len(data))
	for jsonKey, jsonValue := range data {
		key := cfg.UnmarshalJsonKey(jsonKey)
		value, err := unmarshalJson(key, jsonValue)
		if err != nil {
			return keyUnmarshalError{key, err}
		}
		paramsMap[key] = value
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

func unmarshalJson(key string, data []byte) (any, error) {
	unmarshller, found := cfg.UnmarshalJsonParam[key]
	if found {
		return unmarshller(data)
	}
	var value any
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
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

func marshalXml(key string, en *xml.Encoder, value any) error {
	marshaller, found := cfg.MarshalXmlParam[key]
	if found {
		return marshaller(en, xml.StartElement{Name: xml.Name{Local: cfg.MarshalXMLKey(key)}}, value)
	}
	return en.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: cfg.MarshalXMLKey(key)}})
}

var _ xml.Unmarshaler = (*Error)(nil)

func (e *Error) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	var code Code
	var codeFound bool
	var message string
	var messageFound bool
	var cause error
	var paramsMap = map[string]any{}
	for {
		token, _ := d.Token()
		if token == nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		key := cfg.UnmarshalXMLKey(start.Name.Local)
		switch key {
		case keyCode:
			codeFound = true
			codeValue, err := unmarshalXml(keyCode, d, start)
			if err != nil {
				return keyUnmarshalError{keyCode, err}
			}
			code, ok = codeValue.(Code)
			if !ok {
				return keyCastError{keyCode}
			}
		case keyMessage:
			messageFound = true
			messageValue, err := unmarshalXml(keyMessage, d, start)
			if err != nil {
				return keyUnmarshalError{keyMessage, err}
			}
			message, ok = messageValue.(string)
			if !ok {
				return keyCastError{keyMessage}
			}
		case keyCause:
			causeValue, err := unmarshalXml(keyCause, d, start)
			if err != nil {
				return keyUnmarshalError{keyCause, err}
			}
			cause, ok = causeValue.(error)
			if !ok {
				return keyCastError{keyCause}
			}
		case keyStackTrace:
			break
		default:
			value, err := unmarshalXml(key, d, start)
			if err != nil {
				return keyUnmarshalError{key, err}
			}
			paramsMap[key] = value
		}
	}
	if !codeFound {
		return keyMissingError{keyCode}
	}
	if !messageFound {
		return keyMissingError{keyMessage}
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

func unmarshalXml(key string, d *xml.Decoder, start xml.StartElement) (any, error) {
	unmarshaller, found := cfg.UnmarshalXmlParam[key]
	if found {
		return unmarshaller(d, start)
	}
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return nil, err
	}
	return value, nil
}

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

type keyMissingError struct {
	key string
}

func (e keyMissingError) Error() string {
	return fmt.Sprintf("missing %s", e.key)
}

type keyCastError struct {
	key string
}

func (e keyCastError) Error() string {
	return fmt.Sprintf("cast %s unmarshaled value", e.key)
}
