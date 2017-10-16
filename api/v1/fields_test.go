// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"encoding/json"
	"testing"

	"github.com/praelatus/praelatus/models"
)

func fieldSchemeFromJSON(jsn []byte) (interface{}, error) {
	var tk models.FieldScheme
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func fieldSchemesFromJSON(jsn []byte) (interface{}, error) {
	var tk []models.FieldScheme
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func toFieldSchemes(v interface{}) []models.FieldScheme {
	return v.([]models.FieldScheme)
}

func toFieldScheme(v interface{}) models.FieldScheme {
	return v.(models.FieldScheme)
}

var fieldSchemeRouteTests = []routeTest{
	{
		Name:      "Get FieldScheme",
		Admin:     true,
		Endpoint:  "/api/v1/fieldschemes/59e3f2026791c08e74da1bb2",
		Converter: fieldSchemeFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toFieldScheme(v)

			if fs.ID != "59e3f2026791c08e74da1bb2" {
				t.Errorf("Expected 59e3f2026791c08e74da1bb2 Got: %s", fs.ID)
			}
		},
	},

	{
		Name:      "Get All FieldSchemes",
		Admin:     true,
		Endpoint:  "/api/v1/fieldschemes",
		Converter: fieldSchemesFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toFieldSchemes(v)

			if len(fs) != 1 {
				t.Errorf("Expected 1 FieldScheme got %d", len(fs))
			}

			if fs[0].ID == "" {
				t.Errorf("Expected an ID Got None")
			}
		},
	},

	{
		Name:      "Create FieldScheme",
		Admin:     true,
		Method:    "POST",
		Endpoint:  "/api/v1/fieldschemes",
		Converter: fieldSchemeFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toFieldScheme(v)

			if fs.ID == "" {
				t.Errorf("Expected An ID but got None")
			}
		},
	},

	{
		Name:     "Remove FieldScheme",
		Endpoint: "/api/v1/fieldschemes/59e3f2026791c08e74da1bb2",
		Method:   "DELETE",
	},
}
