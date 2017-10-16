// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"encoding/json"
	"testing"

	"github.com/praelatus/praelatus/models"
)

func workflowFromJSON(jsn []byte) (interface{}, error) {
	var tk models.Workflow
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func workflowsFromJSON(jsn []byte) (interface{}, error) {
	var tk []models.Workflow
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func toWorkflows(v interface{}) []models.Workflow {
	return v.([]models.Workflow)
}

func toWorkflow(v interface{}) models.Workflow {
	return v.(models.Workflow)
}

var workflowRouteTests = []routeTest{
	{
		Name:      "Get Workflow",
		Admin:     true,
		Endpoint:  "/api/v1/workflows/59e3f2026791c08e74da1bb2",
		Converter: workflowFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toWorkflow(v)

			if fs.ID != "59e3f2026791c08e74da1bb2" {
				t.Errorf("Expected 59e3f2026791c08e74da1bb2 Got: %s", fs.ID)
			}
		},
	},

	{
		Name:      "Get All Workflows",
		Admin:     true,
		Endpoint:  "/api/v1/workflows",
		Converter: workflowsFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toWorkflows(v)

			if len(fs) != 1 {
				t.Errorf("Expected 1 Workflow got %d", len(fs))
			}

			if fs[0].ID == "" {
				t.Errorf("Expected an ID Got None")
			}
		},
	},

	{
		Name:      "Create Workflow",
		Admin:     true,
		Method:    "POST",
		Endpoint:  "/api/v1/workflows",
		Converter: workflowFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			fs := toWorkflow(v)

			if fs.ID == "" {
				t.Errorf("Expected An ID but got None")
			}
		},
	},

	{
		Name:     "Remove Workflow",
		Endpoint: "/api/v1/workflows/59e3f2026791c08e74da1bb2",
		Method:   "DELETE",
	},
}
