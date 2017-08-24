package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
)

func workflowRouter(router *mux.Router) {
	router.HandleFunc("/workflows", getAllWorkflows).Methods("GET")
	router.HandleFunc("/workflows", createWorkflow).Methods("POST")
	router.HandleFunc("/workflows/{id}", singleWorkflow)
}

func createWorkflow(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var workflow models.Workflow

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&workflow)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	workflow.ID = bson.NewObjectId()

	err = getCollection(config.WorkflowCollection).Insert(workflow)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, workflow)
}

func getAllWorkflows(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var ws []models.Workflow

	var query bson.M
	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
		query = bson.M{"name": bson.M{"$regex": q, "$options": "i"}}
	}

	err := getCollection(config.WorkflowCollection).Find(query).All(&ws)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.APIError(err.Error()))
		return
	}

	utils.SendJSON(w, ws)
}

func singleWorkflow(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var f models.Workflow
	id := bson.ObjectIdHex(mux.Vars(r)["id"])
	coll := getCollection(config.WorkflowCollection)

	var err error

	switch r.Method {
	case "GET":
		err = coll.FindId(id).One(&f)
	case "DELETE":
		err = coll.RemoveId(id)
	case "PUT":
		var workflow models.Workflow

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&workflow)
		if err != nil {
			break
		}

		err = coll.UpdateId(id, &workflow)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.APIMsg(err.Error()))
		return
	}

	if f.Name != "" {
		utils.SendJSON(w, f)
		return
	}

	utils.SendJSON(w, map[string]string{})
}
