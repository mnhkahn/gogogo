// Package func_to_handler ...
package func_to_handler

import (
	"encoding/json"
	"fmt"
	"github.com/mnhkahn/gogogo/xreflect"
	"net/http"
	"reflect"
)

type FuncToHandler struct {
	paramIns  []reflect.Type
	paramOuts []reflect.Type
	method    reflect.Value
}

func NewFuncToHandler(fn interface{}) *FuncToHandler {
	f := new(FuncToHandler)
	//rt := reflect.TypeOf(fn)
	rv := reflect.ValueOf(fn)

	for i := 0; i < rv.Type().NumIn(); i++ {
		t := rv.Type().In(i)
		f.paramIns = append(f.paramIns, t)
	}
	for i := 0; i < rv.Type().NumOut(); i++ {
		t := rv.Type().Out(i)
		f.paramOuts = append(f.paramOuts, t)
	}
	f.method = rv
	return f
}

func (f *FuncToHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get parameter in
	query := r.URL.Query()
	var inValues []reflect.Value
	for i, in := range f.paramIns {
		value := query.Get(fmt.Sprintf("in%d", i))

		v, err := xreflect.StringToReflectValue(value, in.Kind().String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		inValues = append(inValues, v)
	}

	// call function
	outValues := f.method.Call(inValues)

	res := make(map[string]interface{}, len(outValues))
	for i, o := range outValues {
		inter := o.Interface()
		if o.Kind() == reflect.Slice {
			if o.Type().Elem().Kind() == reflect.Uint8 {
				inter = string(o.Interface().([]byte))
			}
		}
		res[fmt.Sprintf("out%d", i)] = inter
	}

	buf, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(buf)
}
