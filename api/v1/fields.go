package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
)

func fieldRouter(router *mux.Router) {
	router.HandleFunc("/fieldschemes", getAllFieldSchemes).Methods("GET")
	router.HandleFunc("/fieldschemes", createFieldScheme).Methods("POST")
	router.HandleFunc("/fieldschemes/{id}", singleFieldScheme)
}

func createFieldScheme(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var f map[string]models.FieldScheme

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&f)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	fs := f["fieldScheme"]
	fs.ID = bson.NewObjectId()

	err = getCollection(config.FieldSchemeCollection).Insert(fs)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, fs)
}

func getAllFieldSchemes(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var fs []models.FieldScheme

	var query bson.M
	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
		query = bson.M{"name": bson.M{"$regex": q, "$options": "i"}}
	}

	fmt.Println("query", query)

	err := getCollection(config.FieldSchemeCollection).Find(query).All(&fs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.APIError(err.Error()))
		return
	}

	utils.SendJSON(w, fs)
}

func singleFieldScheme(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil || !u.IsAdmin {
		utils.APIErr(w, http.StatusForbidden,
			"you must be logged in as an administrator")
		return
	}

	var f models.FieldScheme
	id := bson.ObjectIdHex(mux.Vars(r)["id"])
	coll := getCollection(config.FieldSchemeCollection)

	var err error

	switch r.Method {
	case "GET":
		err = coll.FindId(id).One(&f)
	case "DELETE":
		err = coll.RemoveId(id)
	case "PUT":
		var jr map[string]models.FieldScheme

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&jr)
		if err != nil {
			break
		}

		f, ok := jr["fieldScheme"]
		if !ok {
			err = errors.New("invalid object schema")
			break
		}

		err = coll.UpdateId(id, &f)
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
