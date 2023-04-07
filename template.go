package errors

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Code string

type Params = map[string]any

var MessageParamNotFoundHandler = func(template Template, param string, trace []string) {}

type Error struct {
	template Template
	params   map[string]any
	trace    []string
}

func newError(template Template, params Params) Error {
	trace := stackTrace(3)
	templateParams := template.messageParams()
	for _, param := range templateParams {
		_, found := params[param]
		if !found {
			MessageParamNotFoundHandler(template, param, trace)
		}
	}
	return Error{
		template: template,
		params:   params,
		trace:    trace,
	}
}

func (e Error) Code() Code {
	return e.template.Code
}

func (e Error) Message() string {
	return e.template.Format(e.params)
}

func (e Error) Error() string {
	return e.Message()
}

func (e Error) Trace() []string {
	return e.trace
}

type Template struct {
	Code    Code
	Message string
	Params  Params
}

func (t Template) Format(params Params) string {
	return format(t.Message, t.mergeParams(params))
}

func (t Template) messageParams() []string {
	return params(t.Message)
}

func (t Template) mergeParams(params Params) Params {
	return mergeParams(t.Params, params)
}

func mergeParams(params1 Params, params2 Params) Params {
	mergedParams := make(Params, len(params1)+len(params2))
	for paramId, param := range params1 {
		mergedParams[paramId] = param
	}
	for paramId, param := range params2 {
		mergedParams[paramId] = param
	}
	return mergedParams
}

func params(str string) []string {
	var result []string
	tokens := tokenize(str, brackets)
	index := 0
	for index < len(tokens) {
		var param string
		var ok bool
		param, ok, index = readId(tokens, index)
		if ok {
			result = append(result, param)
		} else {
			_, _, index = readToken(tokens, index)
		}
	}
	return result
}

func format(str string, params map[string]any) string {
	tokens := tokenize(str, brackets)
	sb := strings.Builder{}
	index := 0
	for index < len(tokens) {
		var paramId string
		var ok bool
		var data string
		paramId, ok, index = readId(tokens, index)
		if ok {
			param, found := params[paramId]
			if found {
				data = fmt.Sprintf("%v", param)
			} else {
				data = fmt.Sprintf("{%s}", paramId)
			}
		} else {
			data, _, index = readToken(tokens, index)
		}
		sb.WriteString(data)
	}
	return sb.String()
}

func readId(tokens []string, startIndex int) (string, bool, int) {
	index := startIndex
	if index >= len(tokens) {
		return "", false, startIndex
	}
	if tokens[index] != openBracket {
		return "", false, startIndex
	}
	index++
	for {
		if index >= len(tokens) {
			return "", false, startIndex
		}
		token := tokens[index]
		if token == closeBracket {
			return strings.Join(tokens[startIndex+1:index], ""), true, index + 1
		}
		index++
	}
}

func readToken(tokens []string, index int) (string, bool, int) {
	if index >= len(tokens) {
		return "", false, index
	}
	return tokens[index], false, index + 1
}

const (
	openBracket  = "ob"
	closeBracket = "cb"
)

var brackets = map[rune]string{
	'{': openBracket,
	'}': closeBracket,
}

func tokenize(str string, escapeeMap map[rune]string) []string {
	tokens := make([]string, 0, utf8.RuneCountInString(str))
	index := 0
	escapee := false
	for index < len(str) {
		r, n := utf8.DecodeRuneInString(str[index:])
		switch {
		case !escapee && r == '\\':
			escapee = true
		case escapee:
			tokens = append(tokens, str[index:index+n])
			escapee = false
		default:
			token, found := escapeeMap[r]
			if found {
				tokens = append(tokens, token)
			} else {
				tokens = append(tokens, str[index:index+n])
			}
		}
		index += n
	}
	return tokens
}

var NotFound = Template{
	Code:    "NOT_FOUND",
	Message: "{target} not found",
}

var AlreadyExists = Template{
	Code:    "ALREADY_EXISTS",
	Message: "{target} already exists",
}

var NotAllowed = Template{
	Code:    "NOT_ALLOWED",
	Message: "Not allowed: {cause}",
}

var Timeout = Template{
	Code:    "TIMEOUT",
	Message: "Timeout",
}

var Canceled = Template{
	Code:    "CANCELED",
	Message: "Operation canceled",
}

var NotImplemented = Template{
	Code:    "NOT_IMPLEMENTED",
	Message: "Not implemented",
}

var InternalError = Template{
	Code:    "INTERNAL_ERROR",
	Message: "Internal error",
}
