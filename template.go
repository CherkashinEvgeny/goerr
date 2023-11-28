package errors

import (
	"strings"
	"text/template"
)

type Template struct {
	Code    Code
	Message func(params map[string]any) string
	Params  Params
}

func Message(str string) func(params map[string]any) string {
	tmpt, err := template.New("").Funcs(template.FuncMap{
		"default": func(defaultVal any, val any) any {
			if val == nil {
				return defaultVal
			}
			return val
		},
	}).Parse(str)
	if err != nil {
		panic(err)
	}
	return func(params map[string]any) string {
		sb := strings.Builder{}
		err = tmpt.Execute(&sb, params)
		if err != nil {
			return str
		}
		return sb.String()
	}
}
