package errors

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var basePath string

func init() {
	d, _ := os.Getwd()
	if path.Base(d) == "bin" {
		d = path.Dir(d)
	}
	basePath = path.Dir(d)
}

func stackTrace(index int) []string {
	steps := make([]string, 0, 10)
	for {
		t := stackTraceItem(index)
		if t == "" {
			break
		}
		steps = append(steps, t)
		index++
	}
	return steps
}

func stackTraceItem(index int) string {
	pc, filePath, line, ok := runtime.Caller(index)
	if line == 0 || !strings.Contains(filePath, basePath) {
		return ""
	}
	item := strings.ReplaceAll(filePath, basePath, "") + ":" + strconv.Itoa(line)
	if !ok {
		return item
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return item
	}
	item += " " + path.Base(fn.Name())
	return item
}
