// Package func_to_handler ...
package func_to_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/mnhkahn/gogogo/app"
	"github.com/mnhkahn/gogogo/xreflect"
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

func (f *FuncToHandler) ServeFunc(c *app.Context) error {
	if app.DefaultHandler.RecoverFunc != nil {
		defer app.DefaultHandler.RecoverFunc(c)
	}
	f.ServeHTTP(c.ResponseWriter, c.Request)
	return nil
}

func (f *FuncToHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get parameter in
	query := r.URL.Query()
	var inValues []reflect.Value
	for i, in := range f.paramIns {
		key := fmt.Sprintf("in%d", i)
		value := ""

		// if request is post, get value from body
		if r.Method == http.MethodPost {
			body := make(map[string]string, len(f.paramIns))
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			value = body[key]
		} else {
			// if request is get, get value from query
			value = query.Get(key)
		}

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
