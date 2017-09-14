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

func fieldRouter(router *mux.Router) {
	router.HandleFunc("/fieldschemes", getAllFieldSchemes).Methods("GET")
	router.HandleFunc("/fieldschemes", createFieldScheme).Methods("POST")
	router.HandleFunc("/fieldschemes/{id}", singleFieldScheme)
}

func createFieldScheme(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	var fs models.FieldScheme

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&fs)
	if err != nil {
		utils.Error(w, err)
		return
	}

	fs, err = Repo.Fields().Create(u, fs)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, fs)
}

func getAllFieldSchemes(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)

	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
	}

	fs, err := Repo.Fields().Search(u, q)
	if err != nil {
		utils.Error(w, err)
		return
	}

	utils.SendJSON(w, fs)
}

func singleFieldScheme(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)

	id := mux.Vars(r)["id"]

	var f models.FieldScheme
	var err error

	switch r.Method {
	case "GET":
		f, err = Repo.Fields().Get(u, id)
	case "DELETE":
		err = Repo.Fields().Delete(u, id)
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&f)
		if err != nil {
			break
		}

		err = Repo.Fields().Update(u, id, f)
	}

	if err != nil {
		utils.Error(w, err)
		return
	}

	if f.Name != "" {
		utils.SendJSON(w, f)
		return
	}

	utils.SendJSON(w, map[string]string{})
}
