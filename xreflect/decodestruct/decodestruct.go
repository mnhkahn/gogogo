// Package decodestruct
package decodestruct

import (
	"reflect"
	"strings"
)

// FuncMap is the name with func.
type FuncMap map[string]func() (interface{}, error)

// Decode dto by tagName. It will auto call func in FuncMap by tag name and will auto set value by dto's field.
func Decode(dto interface{}, tagName string, funcMap FuncMap) (interface{}, error) {
	rv := reflect.ValueOf(dto).Elem()
	rt := reflect.TypeOf(dto).Elem()
	for i := 0; i < rv.NumField(); i++ {
		fieldName := rt.Field(i).Name
		// if it has a tag, use tag instead.
		if tagName != "" {
			if tag := strings.TrimSpace(rt.Field(i).Tag.Get(tagName)); tag != "" {
				fieldName = tag
			}
		}

		// get func by field name
		var fn func() (interface{}, error)
		var ok bool
		if fn, ok = funcMap[fieldName]; !ok {
			continue
		}

		// call tag function
		res, err := fn()
		if err != nil {
			return dto, err
		}

		// save value
		if res != nil {
			rv.Field(i).Set(reflect.ValueOf(res))
		}
	}
	return nil, nil
}
