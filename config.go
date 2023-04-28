package errors

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type Config struct {
	CollectStackTrace bool

	IsPrivateParam func(name string) bool

	ParamNameToJsonKey   func(name string) string
	MarshalJsonParam     map[string]func(v any) ([]byte, error)
	MarshalJson          func(v any) ([]byte, error)
	ParamNameFromJsonKey func(name string) string
	UnmarshalJsonParam   map[string]func(data []byte, v any) error
	UnmarshalJson        func(data []byte, v any) error

	ParamNameToXMLKey   func(name string) string
	MarshalXmlParam     map[string]func(en *xml.Encoder, start xml.StartElement, v any) error
	MarshalXml          func(en *xml.Encoder, start xml.StartElement, v any) error
	ParamNameFromXMLKey func(name string) string
	UnmarshalXmlParam   map[string]func(d *xml.Decoder, start xml.StartElement, v any) error
	UnmarshalXml        func(d *xml.Decoder, start xml.StartElement, v any) error

	MarshalCause      bool
	MarshalStackTrace bool
}

var cfg = Config{
	CollectStackTrace: true,

	IsPrivateParam: func(name string) bool {
		r, _ := utf8.DecodeRuneInString(name)
		return unicode.IsLower(r)
	},

	ParamNameToJsonKey: func(name string) string {
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
	MarshalJson: json.Marshal,
	ParamNameFromJsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	UnmarshalJsonParam: map[string]func(data []byte, v any) error{
		keyCode: func(data []byte, v any) error {
			return json.Unmarshal(data, &v)
		},
		keyMessage: func(data []byte, v any) error {
			return json.Unmarshal(data, &v)
		},
		keyCause: func(data []byte, v any) error {
			var e Error
			err := json.Unmarshal(data, &e)
			if err == nil {
				reflect.ValueOf(v).Elem().Set(reflect.ValueOf(e))
				return nil
			}
			var str string
			err = json.Unmarshal(data, &str)
			if err != nil {
				reflect.ValueOf(v).Elem().Set(reflect.ValueOf(errors.New(str)))
				return nil
			}
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(errors.New(string(data))))
			return nil
		},
	},
	UnmarshalJson: json.Unmarshal,

	ParamNameToXMLKey: func(name string) string {
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
	},
	MarshalXml: func(en *xml.Encoder, start xml.StartElement, v any) error {
		return en.EncodeElement(v, start)
	},
	ParamNameFromXMLKey: func(name string) string {
		return name
	},
	UnmarshalXmlParam: map[string]func(d *xml.Decoder, start xml.StartElement, v any) error{
		keyCode: func(d *xml.Decoder, start xml.StartElement, v any) error {
			return d.DecodeElement(v, &start)
		},
		keyMessage: func(d *xml.Decoder, start xml.StartElement, v any) error {
			return d.DecodeElement(v, &start)
		},
		keyCause: func(d *xml.Decoder, start xml.StartElement, v any) error {
			var e Error
			err := d.DecodeElement(&e, &start)
			if err == nil {
				reflect.ValueOf(v).Elem().Set(reflect.ValueOf(e))
				return nil
			}
			var str string
			err = d.DecodeElement(&str, &start)
			if err != nil {
				reflect.ValueOf(v).Elem().Set(reflect.ValueOf(errors.New(str)))
				return nil
			}
			return nil
		},
	},
	UnmarshalXml: func(d *xml.Decoder, start xml.StartElement, v any) error {
		return d.DecodeElement(v, &start)
	},

	MarshalCause:      false,
	MarshalStackTrace: false,
}

func Configure(f func(Config) Config) {
	cfg = f(cfg)
}
