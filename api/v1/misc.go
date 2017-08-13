package v1

// Contains various endpoints that aren't full models in their own right.

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/models"
	"github.com/praelatus/backend/models/permission"
)

func miscRouter(router *mux.Router) {
	router.HandleFunc("/permissions", GetAllPermissionSchemes)
	// router.HandleFunc("/labels", GetAllLabels).Methods("GET")
	// router.HandleFunc("/types", GetAllTypes).Methods("GET")
}

// GetAllPermissionSchemes will return a JSON array of all permissionSchemes from the store.
func GetAllPermissionSchemes(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	utils.SendJSONR(w,
		models.JSONRepr{"permissions": permission.Permissions})
}
