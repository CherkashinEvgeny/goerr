package errors

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Config struct {
	MissingTemplateParamHandler func(template Template, param string)
	IsPrivateParam              func(name string) bool
	ToJsonKey                   func(name string) string
	FromJsonKey                 func(name string) string
	ToXMLKey                    func(name string) string
	FromXMLKey                  func(name string) string
	CollectStackTrace           bool
	MarshalStackTrace           bool
}

var cfg = Config{
	MissingTemplateParamHandler: func(template Template, param string) {
		panic(fmt.Errorf("parameter = \"%s\" is missing", param))
	},
	IsPrivateParam: func(name string) bool {
		r, _ := utf8.DecodeRuneInString(name)
		return unicode.IsLower(r)
	},
	ToJsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsLower(r) {
			return name
		}
		return string(unicode.ToLower(r)) + name[n:]
	},
	FromJsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	ToXMLKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	FromXMLKey: func(name string) string {
		return name
	},
	CollectStackTrace: true,
	MarshalStackTrace: false,
}

func Configure(f func(Config) Config) {
	cfg = f(cfg)
}
