package marshaller

import (
	"encoding/json"
	"testing"
)

type testStruct struct {
	Str    string      `json:"str" groups:"public"`
	Int    int         `json:"int" groups:"internal"`
	Nil    interface{} `json:"nil" groups:"internal,public"`
	Nested *testStruct `json:"nested" groups:"internal,public"`
}

func TestMarshaller_MarshalPublic(t *testing.T) {
	data, authErr := MarshalPublic(&testStruct{
		Str: "str",
		Int: 1,
		Nil: nil,
		Nested: &testStruct{
			Str: "nested",
			Int: 2,
			Nil: nil,
		},
	})
	if nil != authErr {
		t.Error(authErr)
	}
	if nil == data {
		t.Error("Data is nil")
	}
	jsonData, err := json.Marshal(data)
	if nil != err {
		t.Error(err)
	}
	if nil == jsonData {
		t.Error("Json data is nil")
	}
	jsonStr := string(jsonData)
	if "" == jsonStr {
		t.Error("Json data is empty")
	}
	if jsonStr != "{\"nested\":{\"nested\":null,\"nil\":null,\"str\":\"nested\"},\"nil\":null,\"str\":\"str\"}" {
		t.Error("Json data is not as expected")
	}
}

func TestMarshaller_MarshalInternal(t *testing.T) {
	data, authErr := MarshalInternal(&testStruct{
		Str: "str",
		Int: 1,
		Nil: nil,
		Nested: &testStruct{
			Str: "nested",
			Int: 2,
			Nil: nil,
		},
	})
	if nil != authErr {
		t.Error(authErr)
	}
	if nil == data {
		t.Error("Data is nil")
	}
	jsonData, err := json.Marshal(data)
	if nil != err {
		t.Error(err)
	}
	if nil == jsonData {
		t.Error("Json data is nil")
	}
	jsonStr := string(jsonData)
	if "" == jsonStr {
		t.Error("Json data is empty")
	}
	if jsonStr != "{\"int\":1,\"nested\":{\"int\":2,\"nested\":null,\"nil\":null},\"nil\":null}" {
		t.Error("Json data is not as expected")
	}
}
