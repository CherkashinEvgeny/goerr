package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	Configure(func(config Config) Config {
		config.MarshalStackTrace = true
		return config
	})
	e := New(NotFound, Resource("Bucket"))
	errBytes, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(errBytes))
}
