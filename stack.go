package errors

import (
	"io"
	"runtime"
	"strconv"
	"strings"
)

const maxDepth = 32

func callers(skip int) Stack {
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[0:n]
}

type Stack []uintptr

func (s Stack) String() string {
	return s.StackTrace().String()
}

func (s Stack) CompactString() string {
	return s.StackTrace().CompactString()
}

func (s Stack) PrettyString() string {
	return s.StackTrace().PrettyString()
}

func (s Stack) StackTrace() StackTrace {
	st := make([]Frame, 0, len(s))
	for _, p := range s {
		st = append(st, Frame(p))
	}
	return st
}

type StackTrace []Frame

func (st StackTrace) String() string {
	return st.CompactString()
}

func (st StackTrace) CompactString() string {
	sb := strings.Builder{}
	st.compactWrite(&sb)
	return sb.String()
}

func (st StackTrace) compactWrite(w io.Writer) {
	_, _ = io.WriteString(w, "[")
	for i, f := range st {
		if i > 0 {
			_, _ = io.WriteString(w, " ")
		}
		f.compactWrite(w)
	}
	_, _ = io.WriteString(w, "]")
}

func (st StackTrace) PrettyString() string {
	sb := strings.Builder{}
	st.prettyWrite(&sb)
	return sb.String()
}

func (st StackTrace) prettyWrite(w io.Writer) {
	for i, f := range st {
		if i > 0 {
			_, _ = io.WriteString(w, "\n")
		}
		f.prettyWrite(w)
	}
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

func (f Frame) String() string {
	return f.CompactString()
}

func (f Frame) CompactString() string {
	sb := strings.Builder{}
	f.compactWrite(&sb)
	return sb.String()
}

func (f Frame) compactWrite(w io.Writer) {
	f.write(w, " ")
}

func (f Frame) PrettyString() string {
	sb := strings.Builder{}
	f.prettyWrite(&sb)
	return sb.String()
}

func (f Frame) prettyWrite(w io.Writer) {
	f.write(w, "\n\t")
}

func (f Frame) write(w io.Writer, separator string) {
	name := f.Func()
	if name == "" {
		file := f.File()
		if file == "" {
			_, _ = io.WriteString(w, "unknown")
			return
		}
		_, _ = io.WriteString(w, file)
		_, _ = io.WriteString(w, ":")
		_, _ = io.WriteString(w, strconv.Itoa(f.Line()))
		return
	}
	_, _ = io.WriteString(w, name)
	file := f.File()
	if file == "" {
		return
	}
	_, _ = io.WriteString(w, separator)
	_, _ = io.WriteString(w, file)
	_, _ = io.WriteString(w, ":")
	_, _ = io.WriteString(w, strconv.Itoa(f.Line()))
}

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}
