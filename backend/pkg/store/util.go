package store

import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/tidwall/buntdb"
)

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func jsonMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	must(err)
	return unsafeBytesToString(data)
}

func jsonMustUnmarshal(s string, v interface{}) {
	must(json.NewDecoder(strings.NewReader(s)).Decode(v))
}

func jsonUnmarshal(s string, v interface{}) error {
	err := json.NewDecoder(strings.NewReader(s)).Decode(v)
	if err != nil {
		log.Printf("json %v\n%v", err, s)
	}
	return err
}

const timeLayout = "2006-01-02T15:04:05Z"

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format(timeLayout)
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(timeLayout, s)
	must(err)
	return t
}

func txGet(tx *buntdb.Tx, key string) string {
	v, _ := tx.Get(key)
	return v
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
