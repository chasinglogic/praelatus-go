package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models"
)

func projectRouter(router *mux.Router) {
	router.HandleFunc("/projects", GetAllProjects).Methods("GET")
	router.HandleFunc("/projects", CreateProject).Methods("POST")

	router.HandleFunc("/projects/{key}", GetProject).Methods("GET")
	router.HandleFunc("/projects/{key}/tickets", GetAllTicketsByProject).Methods("GET")
	router.HandleFunc("/projects/{key}", RemoveProject).Methods("DELETE")
	router.HandleFunc("/projects/{key}", UpdateProject).Methods("PUT")

	router.HandleFunc("/projects/{key}/fields/{ticketType}", GetFieldsForScreen)

	router.HandleFunc("/projects/{key}/roles", GetRolesForProject)
	router.HandleFunc("/projects/{key}/roles/{roleId}/addUser/{userId}", AddUserToRole)

	router.HandleFunc("/projects/{key}/permissionscheme", GetPermissionScheme).Methods("GET")
	router.HandleFunc("/projects/{key}/permissionscheme/{schemeId}", SetPermissionScheme).Methods("POST")
}

// GetProject will get a project by it's project key
func GetProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{ID: 0}
	}

	vars := mux.Vars(r)
	key := vars["key"]

	p := models.Project{
		Key: key,
	}

	err := Store.Projects().Get(*u, &p)
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
		u = &models.User{ID: 0}
	}

	projects, err := Store.Projects().GetAll(*u)
	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	utils.SendJSON(w, projects)
}

// GetAllTicketsByProject will get all the tickets for a given project
func GetAllTicketsByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkey := vars["key"]

	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{ID: 0}
	}

	p := models.Project{Key: pkey}
	err := Store.Projects().Get(*u, &p)
	if err != nil {
		utils.APIErr(w, 500, err.Error())
	}

	tks, err := Store.Tickets().GetAllByProject(*u, p)
	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError("failed to retrieve tickets from the database"))
		log.Println(err)
		return
	}

	utils.SendJSON(w, tks)
}

// CreateProject will create a project based on the JSON representation sent to
// the API
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project

	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		w.WriteHeader(403)
		w.Write(utils.APIError("you must be logged in as a system administrator to create a project"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(400)
		w.Write(utils.APIError("invalid body"))
		log.Println(err)
		return
	}

	err = Store.Projects().Create(*u, &p)
	if err != nil {
		w.WriteHeader(400)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	utils.SendJSON(w, p)
}

// RemoveProject will remove the project indicated by the key passed in as a
// url parameter
func RemoveProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		w.WriteHeader(403)
		w.Write(utils.APIError("you must be logged in as a system administrator to delete a project"))
		return
	}

	err := Store.Projects().Remove(*u, models.Project{Key: key})
	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	w.Write([]byte{})

}

// UpdateProject will update a project based on the JSON representation sent to
// the API
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project

	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		w.WriteHeader(403)
		w.Write(utils.APIError("you must be logged in as a system administrator to update a project"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		w.WriteHeader(400)
		w.Write(utils.APIError("invalid body"))
		log.Println(err)
		return
	}

	err = Store.Projects().Save(*u, p)
	if err != nil {
		w.WriteHeader(400)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	utils.SendJSON(w, p)
}

// GetFieldsForScreen will return the appropriate fields based on the given
// project key and ticket type
func GetFieldsForScreen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ticketType := vars["ticketType"]

	w.Write([]byte(key + " " + ticketType))
}

// GetRolesForProject will return the system roles with the members
// who are configured for the given project
func GetRolesForProject(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusUnauthorized,
			"you must be logged in as a project administrator")
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]

	var p models.Project
	p.Key = key

	err := Store.Projects().Get(*u, &p)
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	roles, err := Store.Roles().GetForProject(*u, p)
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	utils.SendJSON(w, roles)
}

// AddUserToRole takes the project key finds the project then adds the
// user for userID to the role for roleID on that project
func AddUserToRole(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusUnauthorized,
			"you must be logged in as a project administrator")
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	roleID, err := strconv.Atoi(vars["roleId"])
	userID, uerr := strconv.Atoi(vars["userId"])

	if err != nil || uerr != nil {
		utils.APIErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Project
	p.Key = key

	err = Store.Projects().Get(*u, &p)
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	err = Store.Roles().AddUserToRole(*u, models.User{ID: int64(userID)},
		p, models.Role{ID: int64(roleID)})
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	w.Write(utils.APIMsg("successfully added user to role"))
}

// SetPermissionScheme will handle the incoming request and set the
// permission scheme for the appropriate project
func SetPermissionScheme(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusUnauthorized,
			"you must be logged in as a project administrator")
		return
	}

	vars := mux.Vars(r)

	var p models.Project
	p.Key = vars["key"]

	id, err := strconv.Atoi(vars["schemeId"])
	if err != nil {
		utils.APIErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = Store.Projects().Get(*u, &p)
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	err = Store.Projects().
		SetPermissionScheme(*u, p, models.PermissionScheme{ID: int64(id)})
	if err != nil {
		utils.APIErr(w, utils.GetErrorCode(err), err.Error())
		return
	}

	w.Write(utils.APIMsg("successfully set permission scheme"))
}
