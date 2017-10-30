// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package utils contains utility functions used throughout the api
// package
package utils

// Message is a general purpose json struct used primarily for error responses.
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

// APIMessage is a general purpose struct for sending messages to the client,
// generally used for errors
type APIMessage struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// APIMsg is a convenience function for generating an API Message
func APIMsg(msg string, fields ...string) []byte {
	e := APIMessage{
		Message: msg,
	}

	if fields != nil {
		e.Field = strings.Join(fields, ",")
	}

	byt, _ := json.Marshal(e)
	return byt
}

// Success returns the default success message
func Success() []byte {
	return APIMsg("operation completed successfully")
}

// Error will get the appropriate error code and message based on err
func Error(w http.ResponseWriter, err error) {
	code := GetErrorCode(err)
	switch err {
	case repo.ErrUnauthorized:
		APIErr(w, code, http.StatusText(code))
	case repo.ErrNotFound:
		APIErr(w, code, http.StatusText(code))
	default:
		APIErr(w, code, err.Error())
	}
}

// APIErr will send the error message and status code to the
// given ResponseWriter
func APIErr(w http.ResponseWriter, status int, msg string) {
	if status >= 500 {
		log.Println("[ISE] ERROR:", msg)
	}

	w.WriteHeader(status)
	w.Write(APIMsg(msg))
}

// GetErrorCode returns the appropriate http status code for the given
// error
func GetErrorCode(e error) int {
	switch e {
	case repo.ErrUnauthorized:
		return http.StatusUnauthorized
	case repo.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// SendJSON is a convenience function for sending JSON to the given
// ResponseWriter. If v is a models.Sanitizer SendJSON will call
// v.Sanitize() before serializing to JSON.
func SendJSON(w http.ResponseWriter, v interface{}) {
	toJsn := v

	if s, ok := v.(models.Sanitizer); ok {
		toJsn = s.Sanitize()
	}

	resp, err := json.Marshal(toJsn)
	if err != nil {
		w.WriteHeader(500)
		w.Write(APIMsg("Failed to marshal database response to JSON."))
		log.Println(err)
		return
	}

	if resp == nil || string(resp) == "null" {
		w.WriteHeader(404)
		w.Write(APIMsg("not found"))
		return
	}

	w.Write(resp)
}

const requireTag = "required"

// ValidateModel will iterate the struct fields checking the tags for required
// fields. This is used during model creation to validate required data is sent
func ValidateModel(model interface{}) error {
	v := reflect.ValueOf(model)

	for i := 0; i < v.NumField(); i++ {
		tag := strings.ToLower(v.Type().Field(i).Tag.Get(requireTag))

		// Skip if not set to true
		if !strings.Contains(tag, "true") {
			continue
		}

		fe := fmt.Errorf("%s is a required field", v.Type().Field(i).Name)
		val := v.Field(i)
		if !val.IsValid() {
			return fe
		}

		k := val.Kind()
		switch {
		case k >= reflect.Int && k <= reflect.Int64:
			intVal := val.Int()
			if intVal == 0 {
				return fe
			}
		case k >= reflect.Uint64 && k <= reflect.Uint64:
			uintVal := val.Uint()
			if uintVal == 0 {
				return fe
			}
		case k == reflect.Float32 || k == reflect.Float64:
			floatVal := val.Float()
			if floatVal == 0.0 {
				return fe
			}
		case k == reflect.Struct:
			e := ValidateModel(val.Interface())
			if e != nil {
				return e
			}
		case k == reflect.String || k == reflect.Slice || k == reflect.Array || k == reflect.Map:
			if val.Len() == 0 {
				return fe
			}
		}
	}

	return nil
}
