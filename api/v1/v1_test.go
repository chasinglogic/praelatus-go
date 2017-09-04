package v1_test

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/v1"
	"github.com/praelatus/backend/models"
	mgo "gopkg.in/mgo.v2"
)

var router *mux.Router

func init() {
	v1.Conn, _ = mgo.Dial("mongodb://localhost/praelatus")
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
