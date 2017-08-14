package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
	"github.com/praelatus/backend/models/permission"
)

func projectRouter(router *mux.Router) {
	router.HandleFunc("/projects", GetAllProjects).Methods("GET")
	// router.HandleFunc("/projects", CreateProject).Methods("POST")

	router.HandleFunc("/projects/{key}", SingleProject)
}

// SingleProject will get a project by it's project key
func SingleProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var p models.Project
	var err error

	key := mux.Vars(r)["key"]
	coll := getCollection(config.ProjectCollection)

	switch r.Method {
	case "GET":
		err = coll.FindId(key).One(&p)
	case "DELETE":
		err = coll.RemoveId(key)
	case "PUT":
		var jr map[string]models.Project

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&jr)
		if err != nil {
			break
		}

		p, ok := jr["project"]
		if !ok {
			err = errors.New("invalid object schema")
			break
		}

		err = coll.UpdateId(key, &p)
	}

	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	utils.SendJSON(w, p)
}

// GetAllProjects will get all the projects on this instance that the user has
// permissions to
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var projects []models.Project
	query := bson.M{
		"$or": []bson.M{
			{
				"key": bson.M{
					"$in": u.ProjectsMemberOf(),
				},
			},
			{
				"permissions.Anonymous": permission.ViewProject,
			},
		},
	}

	q := r.FormValue("q")
	if q != "" {
		query = bson.M{
			"$and": []bson.M{
				query,
				{
					"$or": []bson.M{
						{
							"name": bson.M{
								"$regex":   q,
								"$options": "i",
							},
						},
						{
							"key": bson.M{
								"$regex":   q,
								"$options": "i",
							},
						},
						{
							"lead": bson.M{
								"$regex":   q,
								"$options": "i",
							},
						},
					},
				},
			},
		}
	}

	err := getCollection(config.ProjectCollection).Find(query).All(&projects)
	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	utils.SendJSON(w, projects)
}