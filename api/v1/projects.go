// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

func projectRouter(router *mux.Router) {
	router.HandleFunc("/projects", getAllProjects).Methods("GET")
	router.HandleFunc("/projects", createProject).Methods("POST")

	router.HandleFunc("/projects/{key}", singleProject)
	router.HandleFunc("/projects/{key}/notifications", getProjectNotifications)
}

// singleProject will get a project by it's project key
func singleProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)

	var p models.Project
	var err error

	key := mux.Vars(r)["key"]

	switch r.Method {
	case "GET":
		p, err = Repo.Projects().Get(u, key)
	case "DELETE":
		err = Repo.Projects().Delete(u, key)
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&p)
		if err != nil {
			break
		}

		err = Repo.Projects().Update(u, key, p)
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, p)
}

// getAllProjects will get all the projects on this instance that the user has
// permissions to
func getAllProjects(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
	}

	projects, err := Repo.Projects().Search(u, q)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, projects)
}

// createProject will create a project in the database. This route requires
// admin.
func createProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.Error(w, repo.ErrUnauthorized)
		return
	}

	var p models.Project

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.ValidateModel(p); err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	p, err = Repo.Projects().Create(u, p)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, p)
}

func getProjectNotifications(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	key := mux.Vars(r)["key"]

	unread := r.FormValue("unread") == "true" || r.FormValue("unread") == "TRUE"

	last, _ := strconv.Atoi(r.FormValue("last"))
	// Last cannot be passed to us as 0 if it is 0 that means either nothing
	// or a non-number was passed to set to the default value of 10
	if last == 0 {
		last = 10
	}

	notifications, err := Repo.Notifications().ForProject(u,
		models.Project{Key: key}, unread, last)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, notifications)
}
