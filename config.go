package errors

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Config struct {
	MissingTemplateParamHandler func(template Template, param string)
	IsPrivateParam              func(name string) bool
	JsonKey                     func(name string) string
	XMLKey                      func(name string) string
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
	JsonKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsLower(r) {
			return name
		}
		return string(unicode.ToLower(r)) + name[n:]
	},
	XMLKey: func(name string) string {
		r, n := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			return name
		}
		return string(unicode.ToUpper(r)) + name[n:]
	},
	CollectStackTrace: true,
	MarshalStackTrace: false,
}

func Configure(f func(Config) Config) {
	cfg = f(cfg)
}
