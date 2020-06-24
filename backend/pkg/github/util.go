package github

import (
	"bytes"
	"encoding/json"
	"unsafe"

	"github.com/ng-vu/ujson"
)

type M map[string]interface{}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func jsonMarshalIndent(v interface{}) string {
	resp, err := json.MarshalIndent(v, "", "  ")
	must(err)
	return unsafeBytesToString(bytes.TrimSpace(resp))
}

func jsonReindent(input []byte, prefix, indentStr string) string {
	b := make([]byte, 0, len(input)*3/2)
	err := ujson.Walk(input, func(indent int, key, value []byte) bool {
		if len(b) != 0 && ujson.ShouldAddComma(value, b[len(b)-1]) {
			b = append(b, ',')
		}
		b = append(b, '\n')
		b = append(b, prefix...)
		for i := 0; i < indent; i++ {
			b = append(b, indentStr...)
		}
		if len(key) > 0 {
			b = append(b, key...)
			b = append(b, `: `...)
		}
		b = append(b, value...)
		return true
	})
	if err != nil {
		panic(err)
	}
	return unsafeBytesToString(bytes.TrimSpace(b))
}
