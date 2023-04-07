package errors

import (
	"fmt"
	"testing"
)

var t = map[rune]string{}

const str = "qwertyuiopasdfghjklzxcvbnm"

func BenchmarkTokenize(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tokenize(str, t)
	}
}

func TestFormat(t *testing.T) {
	hehe := NotFound.Format(Params{
		"target": "test",
	})
	fmt.Println(hehe)
}
