// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1

// Contains various endpoints that aren't full models in their own right.

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models/permission"
)

func miscRouter(router *mux.Router) {
	router.HandleFunc("/permissions", getAllPermissions)
	router.HandleFunc("/system/info", getSysInfo)
	// router.HandleFunc("/labels", GetAllLabels).Methods("GET")
	// router.HandleFunc("/types", GetAllTypes).Methods("GET")
}

// GetAllPermissionSchemes will return a JSON array of all permissionSchemes from the store.
func getAllPermissions(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	utils.SendJSON(w, permission.ListOfPermissions)
}

func getSysInfo(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u != nil && u.IsAdmin {
		utils.SendJSON(w, config.Cfg)
	}

	utils.SendJSON(w, config.Cfg.Public())
}

func updateSysInfo(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var c config.Config

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.ValidateModel(c)
	if err != nil {
		utils.Error(w, err)
		return
	}

	c.DBURL = config.Cfg.DBURL
	c.DBName = config.Cfg.DBName
	c.Port = config.Cfg.Port
	c.DataDirectory = config.Cfg.DataDirectory

	config.Cfg = c
	config.Cfg.Save()
}
