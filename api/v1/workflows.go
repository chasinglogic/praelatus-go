package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"
)

func workflowRouter(router *mux.Router) {
	router.HandleFunc("/workflows", getAllWorkflows).Methods("GET")
	router.HandleFunc("/workflows", createWorkflow).Methods("POST")
	router.HandleFunc("/workflows/{id}", singleWorkflow)
}

func createWorkflow(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)

	var workflow models.Workflow

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&workflow)
	if err != nil {
		utils.Error(w, err)
		return
	}

	workflow, err = Repo.Workflows().Create(u, workflow)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, workflow)
}

func getAllWorkflows(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
	}

	ws, err := Repo.Workflows().Search(u, q)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, ws)
}

func singleWorkflow(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	id := mux.Vars(r)["id"]

	var workflow models.Workflow
	var err error

	switch r.Method {
	case "GET":
		workflow, err = Repo.Workflows().Get(u, id)
	case "DELETE":
		err = Repo.Workflows().Delete(u, id)
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&workflow)
		if err != nil {
			break
		}

		err = Repo.Workflows().Update(u, id, workflow)
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	if workflow.Name != "" {
		utils.SendJSON(w, workflow)
		return
	}

	w.Write(utils.Success())
}
