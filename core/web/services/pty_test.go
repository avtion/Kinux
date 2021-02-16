package services

import (
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func Test_dataToVal(t *testing.T) {
	const data = `{"op":1,"data":{"id":"1","container":""}}`

	missionRaw := &struct {
		ID        string `json:"id"`
		Container string `json:"container"`
	}{}
	msg := &struct {
		Op   int `json:"op"`
		Data struct {
			ID        uint   `json:"id"`
			Container string `json:"container"`
		} `json:"data"`
	}{}
	jsoniter.Get([]byte(data), "data").ToVal(missionRaw)
	jsoniter.Get([]byte(data)).ToVal(msg)
	t.Log(missionRaw)
	t.Log(msg)
}
