package utils

import (
	"net/http/httptest"
	"testing"
)

type TestType struct {
	Name string
}

func TestGetType(t *testing.T) {
	if getType(TestType{}) != "testType" {
		t.Errorf("Failure, expected test_type. Get %s\n",
			getType(TestType{}))
	}
}

func TestSendJSON(t *testing.T) {
	tt := TestType{"Test"}
	recorder := httptest.NewRecorder()

	SendJSON(recorder, tt)

	if "{\"testType\":{\"Name\":\"Test\"}}" != string(recorder.Body.Bytes()) {
		t.Fail()
	}
}
