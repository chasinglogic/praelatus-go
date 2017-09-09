package v1

// Contains various endpoints that aren't full models in their own right.

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/models/permission"
)

func miscRouter(router *mux.Router) {
	router.HandleFunc("/permissions", getAllPermissions)
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
