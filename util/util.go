package util

import (
	"crypto/rand"
	"encoding/base64"
	"mime"
	"net/http"
	"strings"
)

func RandomStr(i int) string {
	b := make([]byte, i)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) //yes thats right, panic
	}
	return base64.URLEncoding.EncodeToString(b)
}

func is(mediaType string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		hdr := r.Header.Get("Accept")
		for _, value := range strings.Split(hdr, ",") {
			mediaType, _, err := mime.ParseMediaType(value)
			if err != nil {
				continue
			}
			if strings.HasPrefix(mediaType, mediaType) {
				return true
			}
		}
		return false
	}
}

var IsHTML = is("text/html")
var IsJSON = is("application/json")

// func Merge(dst, src interface{}) {
// 	dstVal := reflect.ValueOf(dst).Elem()
// 	srcVal := reflect.ValueOf(src).Elem()

// 	for i := 0; i < dstVal.NumField(); i++ {
// 		dstField := dstVal.Field(i)
// 		srcField := srcVal.Field(i)

// 		if dstField.Kind() == reflect.Ptr && dstField.IsNil() &&
// 			!srcField.IsNil() {
// 			dstField.Set(srcField)
// 		}
// 	}
// }
