// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/v1"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

var router *mux.Router

func init() {
	v1.Repo = repo.NewMockRepo()
	router = api.Routes()
}

func testLogin(w http.ResponseWriter, r *http.Request) {
	u := models.User{
		Username: "foouser",
		Password: "foopass",
		Email:    "foo@foo.com",
		FullName: "Foo McFooserson",
		IsAdmin:  false,
		IsActive: true,
		Settings: models.Settings{},
	}

	err := middleware.SetUserSession(u, w)
	if err != nil {
		panic(err)
	}

	r.Header.Add("Authorization", w.Header().Get("Token"))
}

func testAdminLogin(w http.ResponseWriter, r *http.Request) {
	u := models.User{
		Username: "foouser",
		Password: "foopass",
		Email:    "foo@foo.com",
		FullName: "Foo McFooserson",
		IsAdmin:  true,
		IsActive: true,
		Settings: models.Settings{},
	}

	err := middleware.SetUserSession(u, w)
	if err != nil {
		panic(err)
	}

	r.Header.Add("Authorization", w.Header().Get("Token"))
}

type routeTest struct {
	Name         string
	Endpoint     string
	Method       string
	Body         interface{}
	Admin        bool
	Login        bool
	ExpectedCode int
	Converter    func(jsn []byte) (interface{}, error)
	Validator    func(v interface{}, t *testing.T)
}

func testRoutes(routes []routeTest, t *testing.T) {
	for _, route := range routes {
		var body io.Reader

		if route.Body != nil {
			byt, err := json.Marshal(route.Body)
			if err != nil {
				t.Errorf("[%s] Error Marshalling Body: %s",
					route.Name, err.Error())
			}

			body = bytes.NewReader(byt)
		}

		method := "GET"
		if route.Method != "" {
			method = route.Method
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, route.Endpoint, body)

		if route.Login {
			testLogin(w, r)
		}

		if route.Admin {
			testAdminLogin(w, r)
		}

		router.ServeHTTP(w, r)

		// t.Log(w.Body)

		expectedCode := 200
		if route.ExpectedCode != 0 {
			expectedCode = route.ExpectedCode
		}

		if expectedCode != w.Code {
			t.Errorf("[%s] Expected Status Code: %d Got: %d",
				route.Name, expectedCode, w.Code)
		}

		if route.Validator != nil {
			v, err := route.Converter(w.Body.Bytes())
			if err != nil {
				t.Errorf("[%s] Error Unmarshalling Response: %s", route.Name, err.Error())
			}

			route.Validator(v, t)
		}
	}
}
