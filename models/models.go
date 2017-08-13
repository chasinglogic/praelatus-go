// Package models contains all of our models and utility functions for
// creating and modifying them
package models

import (
	"encoding/json"
	"fmt"
)

// JSONRepr is used to easily format results into the form Ember.js expects.
type JSONRepr map[string]interface{}

func jsonString(i interface{}) string {
	b, e := json.MarshalIndent(i, "", "\t")
	if e != nil {
		fmt.Println(e)
		return ""
	}

	return string(b)
}
