package errors

import (
	"runtime"
	"strconv"
	"strings"
)

const maxDepth = 32

func trace(skip int) StackTrace {
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	st := make([]Frame, 0, n)
	for i := 0; i < n; i++ {
		st = append(st, Frame(pcs[i]))
	}
	return st
}

type StackTrace []Frame

func (s StackTrace) String() string {
	sb := strings.Builder{}
	for index, frame := range s {
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

type Frame uintptr

func (f Frame) Func() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	return fn.Name()
}

func (f Frame) File() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}
