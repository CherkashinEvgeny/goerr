package errors

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var traceAllow string

func init() {
	d, _ := os.Getwd()
	if path.Base(d) == "bin" {
		d = path.Dir(d)
	}
	traceAllow = path.Dir(d)
}

func Traces(step int) []string {
	tr := make([]string, 0, 10)
	for i := step; true; i++ {
		t := trace(i)
		if t == "" {
			break
		}
		tr = append(tr, t)
	}

	return tr
}

func trace(step int) (kind string) {
	pc, filePath, line, ok := runtime.Caller(step)
	if line == 0 || !strings.Contains(filePath, traceAllow) {
		return
	}
	kind = strings.ReplaceAll(filePath, traceAllow, "") + ":" + strconv.Itoa(line)
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			kind += " " + path.Base(fn.Name())
		}
	}
	return
}
