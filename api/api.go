// Package api has our router and HTTP handlers for all of the available api
// routes
package api

import (
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/config"

	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/api/v1"
)

func routes(router *mux.Router) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs := []string{}

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			rs = append(rs, t)
			return nil
		})

		utils.SendJSON(w, rs)
	})
}

// Routes will return the mux.Router which contains all of the api routes
func Routes() *mux.Router {
	context := config.ContextPath()

	router := mux.NewRouter()
	api := router.PathPrefix(context + "/api").Subrouter()
	v1r := api.PathPrefix("/v1").Subrouter()

	// setup v1 routes
	v1.Routes(v1r)

	// setup latest routes
	v1.Routes(api)

	// setup routes endpoints
	v1r.HandleFunc("/routes", routes(v1r)).Methods("GET")
	api.HandleFunc("/routes", routes(api)).Methods("GET")
	router.HandleFunc("/routes", routes(router)).Methods("GET")

	router.HandleFunc(context+"/",
		func(w http.ResponseWriter, r *http.Request) {
			path := strings.Split(r.URL.Path, "/")
			root := path[1]

			// TODO handle complex context paths (i.e. if
			// we have a context path of /my/praelatus
			// this will not work.)
			if context != "" {
				root = path[2]
			}

			switch root {
			case "api":
				api.ServeHTTP(w, r)
			default:
				router.ServeHTTP(w, r)
			}

			http.ServeFile(w, r, "client/index.html")
		})

	router.PathPrefix(context + "/static/").Handler(
		http.StripPrefix(context+"/static/",
			http.FileServer(http.Dir("client/static/"))))

	return router
}

// New will start running the api on the given port
func New(conn *mgo.Session) http.Handler {
	v1.Conn = conn

	router := Routes()

	middleware.Cache = middleware.NewMongoCache(conn)
	return middleware.LoadMw(router)
}