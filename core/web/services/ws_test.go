package services

import (
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func Test_parseUnknownInterface(t *testing.T) {
	var testStruct = &struct {
		Msg    string `json:"msg"`
		Number int    `json:"number"`
	}{}
	jsoniter.Get([]byte(`{"op":1, "data": {"msg": "ok", "number": 2}}`), "data").ToVal(testStruct)
	t.Log(testStruct)
}
