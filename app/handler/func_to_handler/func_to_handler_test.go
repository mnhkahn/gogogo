// Package func_to_handler
package func_to_handler

import (
	"bytes"
	"encoding/json"
	"github.com/ChimeraCoder/gojson"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestFuncToHandlerAtoi(t *testing.T) {
	fn := NewFuncToHandler(strconv.Atoi)

	req, err := http.NewRequest("GET", "/health-check?in0=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	fn.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"out0":1,"out1":null}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestFuncToHandlerJsonFormat(t *testing.T) {
	fn := NewFuncToHandler(func(data string) ([]byte, error) {
		var out bytes.Buffer
		err := json.Indent(&out, []byte(data), "", "    ")
		if err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	})

	req, err := http.NewRequest("GET", `/health-check?in0={"a":1}`, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	fn.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"out0":"{\n    \"a\": 1\n}","out1":null}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestFuncToHandlerEscape(t *testing.T) {
	fn := NewFuncToHandler(url.QueryEscape)

	req, err := http.NewRequest("GET", "/health-check?in0=a b", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	fn.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"out0":"a+b"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestFuncToHandlerJsonToStruct(t *testing.T) {
	fn := NewFuncToHandler(func(data string) (string, error) {
		var parser gojson.Parser = gojson.ParseJson
		if output, err := gojson.Generate(bytes.NewBufferString(data), parser, "Foo", "main", []string{"json"}, false); err != nil {
			return "", err
		} else {
			return string(output), nil
		}
	})

	req, err := http.NewRequest("GET", `/health-check?in0={"a":1}`, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	fn.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := "{\"out0\":\"package main\\n\\ntype Foo struct {\\n\\tA int64 `json:\\\"a\\\"`\\n}\\n\",\"out1\":null}"
	assert.Equal(t, expected, rr.Body.String())
}
