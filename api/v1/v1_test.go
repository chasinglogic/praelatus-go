// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"net/http"

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
