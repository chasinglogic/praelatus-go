// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"encoding/json"
	"testing"

	"github.com/praelatus/praelatus/models"
)

func userFromJSON(jsn []byte) (interface{}, error) {
	var tk models.User
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func usersFromJSON(jsn []byte) (interface{}, error) {
	var tk []models.User
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func toUsers(v interface{}) []models.User {
	return v.([]models.User)
}

func toUser(v interface{}) models.User {
	return v.(models.User)
}

var userRouteTests = []routeTest{
	{
		Name:      "Get User",
		Converter: userFromJSON,
		Endpoint:  "/api/v1/users/testuser",
		Validator: func(v interface{}, t *testing.T) {
			user := toUser(v)

			if user.Username != "testuser" {
				t.Errorf("Expected testuser Got: %s", user.Username)
			}

			if user.Password != "" {
				t.Error("Sent password back with user.")
			}
		},
	},

	{
		Name:      "Get All Users",
		Converter: usersFromJSON,
		Endpoint:  "/api/v1/users",
		Validator: func(v interface{}, t *testing.T) {
			users := toUsers(v)

			if len(users) <= 1 {
				t.Errorf("Expected more than 1 user got %d", len(users))
				return
			}

			for _, u := range users {
				if u.Password != "" {
					t.Error("Sent password back with user.")
				}
			}
		},
	},

	{
		Name:      "Create User",
		Admin:     true,
		Converter: userFromJSON,
		Method:    "POST",
		Endpoint:  "/api/v1/users",
		Body: models.User{
			Username: "fakeuser",
			Password: "fakepassword",
			Email:    "fake@fake.com",
			FullName: "The Real Faker",
		},
		Validator: func(v interface{}, t *testing.T) {
			user := toUser(v)

			if user.Username != "testuser" {
				t.Errorf("Expected testuser Got: %s", user.Username)
			}

			if user.Password != "" {
				t.Error("Sent password back with user.")
			}
		},
	},

	{
		Name:     "Remove User",
		Endpoint: "/api/v1/users/fakeuser",
		Admin:    true,
		Method:   "DELETE",
	},
}
