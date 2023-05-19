package errors

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Config struct {
	CollectStackTrace bool

	IsPrivateParam func(name string) bool

	MarshalJsonKey     func(name string) string
	MarshalJsonParam   map[string]func(value any) ([]byte, error)
	UnmarshalJsonKey   func(name string) string
	UnmarshalJsonParam map[string]func(data []byte) (any, error)

	MarshalXMLKey     func(name string) string
	MarshalXmlParam   map[string]func(en *xml.Encoder, start xml.StartElement, value any) error
	UnmarshalXMLKey   func(name string) string
	UnmarshalXmlParam map[string]func(d *xml.Decoder, start xml.StartElement) (any, error)

	MarshalCause      bool
	MarshalStackTrace bool
}

var cfg = Config{
	CollectStackTrace: true,

	IsPrivateParam: func(name string) bool {
		r, _ := utf8.DecodeRuneInString(name)
		return unicode.IsLower(r)
	},

	MarshalJsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsLower(r) {
			return name
		}
		return string(unicode.ToLower(r)) + name[n:]
	},
	MarshalJsonParam: map[string]func(v any) ([]byte, error){
		keyCode: func(v any) ([]byte, error) {
			return json.Marshal(v)
		},
		keyMessage: func(v any) ([]byte, error) {
			return json.Marshal(v)
		},
		keyCause: func(v any) ([]byte, error) {
			e, ok := v.(Error)
			if ok {
				return json.Marshal(e)
			}
			err, ok := v.(error)
			if ok {
				return json.Marshal(err.Error())
			}
			return json.Marshal(v)
		},
		keyStackTrace: func(v any) ([]byte, error) {
			st, ok := v.(StackTrace)
			if !ok {
				return json.Marshal(v)
			}
			strs := make([]string, 0, len(st))
			for _, frame := range st {
				strs = append(strs, fmt.Sprintf("%s %s:%d", frame.Func(), frame.File(), frame.Line()))
			}
			return json.Marshal(strs)
		},
	},
	UnmarshalJsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	UnmarshalJsonParam: map[string]func(data []byte) (any, error){
		keyCode: func(data []byte) (any, error) {
			var code Code
			err := json.Unmarshal(data, &code)
			if err != nil {
				return nil, err
			}
			return code, nil
		},
		keyMessage: func(data []byte) (any, error) {
			var message string
			err := json.Unmarshal(data, &message)
			if err != nil {
				return nil, err
			}
			return message, nil
		},
		keyCause: func(data []byte) (any, error) {
			var e Error
			err := json.Unmarshal(data, &e)
			if err == nil {
				return e, nil
			}
			var str string
			err = json.Unmarshal(data, &str)
			if err == nil {
				return errors.New(str), nil
			}
			return errors.New(string(data)), nil
		},
	},

	MarshalXMLKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	MarshalXmlParam: map[string]func(en *xml.Encoder, start xml.StartElement, v any) error{
		keyCode: func(en *xml.Encoder, start xml.StartElement, v any) error {
			return en.EncodeElement(v, start)
		},
		keyMessage: func(en *xml.Encoder, start xml.StartElement, v any) error {
			return en.EncodeElement(v, start)
		},
		keyCause: func(en *xml.Encoder, start xml.StartElement, v any) error {
			e, ok := v.(Error)
			if ok {
				return en.EncodeElement(e, start)
			}
			err, ok := v.(error)
			if ok {
				return en.EncodeElement(err.Error(), start)
			}
			return en.EncodeElement(v, start)
		},
		keyStackTrace: func(en *xml.Encoder, start xml.StartElement, v any) error {
			st, ok := v.(StackTrace)
			if ok {
				return en.EncodeElement(st.String(), start)
			}
			return en.EncodeElement(v, start)
		},
		keyValidationErrors: func(en *xml.Encoder, start xml.StartElement, v any) error {
			errs, ok := v.(map[string]string)
			if !ok {
				return en.EncodeElement(errs, start)
			}
			err := en.EncodeToken(start)
			if err != nil {
				return err
			}
			for key, value := range errs {
				err = en.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}})
				if err != nil {
					return err
				}
			}
			err = en.EncodeToken(start.End())
			if err != nil {
				return err
			}
			return nil
		},
	},
	UnmarshalXMLKey: func(name string) string {
		return name
	},
	UnmarshalXmlParam: map[string]func(d *xml.Decoder, start xml.StartElement) (any, error){
		keyCode: func(d *xml.Decoder, start xml.StartElement) (any, error) {
			var code Code
			err := d.DecodeElement(&code, &start)
			if err != nil {
				return nil, err
			}
			return code, nil
		},
		keyMessage: func(d *xml.Decoder, start xml.StartElement) (any, error) {
			var message string
			err := d.DecodeElement(&message, &start)
			if err != nil {
				return nil, err
			}
			return message, nil
		},
		keyCause: func(d *xml.Decoder, start xml.StartElement) (any, error) {
			var e Error
			err := d.DecodeElement(&e, &start)
			if err == nil {
				return e, nil
			}
			var str string
			err = d.DecodeElement(&str, &start)
			if err != nil {
				return errors.New(str), nil
			}
			return nil, errors.New(fmt.Sprintf("unprocessable error type %s", start.Name.Local))
		},
		keyValidationErrors: func(d *xml.Decoder, _ xml.StartElement) (any, error) {
			errs := map[string]string{}
			for {
				token, _ := d.Token()
				if token == nil {
					break
				}
				start, ok := token.(xml.StartElement)
				if !ok {
					continue
				}
				key := start.Name.Local
				var value string
				err := d.DecodeElement(&value, &start)
				if err != nil {
					return nil, err
				}
				errs[key] = value
			}
			return errs, nil
		},
	},

	MarshalCause:      false,
	MarshalStackTrace: false,
}

func Configure(f func(*Config)) {
	f(&cfg)
}
