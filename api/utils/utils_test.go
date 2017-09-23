package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/praelatus/praelatus/repo"
)

type TestType struct {
	Name string `json:"name"`
}

func TestError(t *testing.T) {
	tests := map[error]int{
		repo.ErrUnauthorized:    http.StatusUnauthorized,
		errors.New("undefined"): http.StatusInternalServerError,
	}

	for err, expectedStatus := range tests {
		recorder := httptest.NewRecorder()
		Error(recorder, err)
		res := recorder.Result()

		if res.StatusCode != expectedStatus {
			t.Errorf("Expected %d Got %d", expectedStatus, res.StatusCode)
		}
	}
}

func TestSendJSON(t *testing.T) {
	tt := TestType{"Test"}
	recorder := httptest.NewRecorder()

	SendJSON(recorder, tt)

	if "{\"name\":\"Test\"}" != string(recorder.Body.Bytes()) {
		t.Fail()
	}
}
