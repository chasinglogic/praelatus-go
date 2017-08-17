// Package utils contains utility functions used throughout the api
// package
package utils

// Message is a general purpose json struct used primarily for error responses.
import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"unicode"
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

// APIError is a legacy function, deprecated should use APIErr or
// APIMsg as appropriate
func APIError(msg string, fields ...string) []byte {
	return APIMsg(msg, fields...)
}

// APIErr will send the appropriate message and status code to the
// given ResponseWriter
func APIErr(w http.ResponseWriter, status int, msg string) {
	if status >= 500 {
		log.Println(msg)
		w.WriteHeader(status)
		w.Write(APIMsg(http.StatusText(status)))
		return
	}

	w.WriteHeader(status)
	w.Write(APIMsg(msg))
}

// GetErrorCode returns the appropriate http status code for the given
// error
func GetErrorCode(e error) int {
	return http.StatusInternalServerError
}

// Use this to properly wrap JSON in a root element.
func getType(v interface{}) string {
	var name string

	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	// lower case first letter since that's what ember expects.
	lower := []rune(name)
	lower[0] = unicode.ToLower(lower[0])

	return string(lower)
}

// SendJSON is a convenience function for sending JSON to the given
// ResponseWriter it will attempt to convert v into a JSONRepr appropriately
// based on the struct name it's only really useful if v is a single record.
// For a result set convert to JSONRepr yourself then use SendJSONR
func SendJSON(w http.ResponseWriter, v interface{}) {
	repr := map[string]interface{}{
		"data": v,
	}

	resp, err := json.Marshal(repr)
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
