package errors

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	cfg.MarshalCause = true
	hehe, err := xml.Marshal(Wrap(
		New(InternalError),
		InternalError,
		WithReason("test"),
	))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hehe))
	e := Error{}
	// err = xml.Unmarshal(hehe, &e)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println(e)
}

func Benchmark(b *testing.B) {
	b.ReportAllocs()
	j := []byte(`{"hehe": ["123123213"]}`)
	s := json.RawMessage{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(j, &s)
	}
}
