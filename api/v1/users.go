package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"

	"github.com/gorilla/mux"
)

func userRouter(router *mux.Router) {
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")

	router.HandleFunc("/users/{username}", singleUser)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	loggedInUser := middleware.GetUserSession(r)
	if loggedInUser == nil {
		loggedInUser = &models.User{}
	}

	var u *models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !loggedInUser.IsAdmin {
		u.IsAdmin = false
	}

	u, err = models.NewUser(u.Username, u.Password, u.FullName, u.Email, u.IsAdmin)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = middleware.SetUserSession(*u, w)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	var tokenResponse struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}

	tokenResponse.Token = w.Header().Get("Token")
	tokenResponse.User = *u

	utils.SendJSON(w, tokenResponse)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	var query bson.M
	q := r.FormValue("q")
	if q != "" {
		q = strings.Replace(q, "*", ".*", -1)
		query = bson.M{
			"$or": []bson.M{
				{"username": bson.M{"$regex": q, "$options": "i"}},
				{"email": bson.M{"$regex": q, "$options": "i"}},
				{"fullname": bson.M{"$regex": q, "$options": "i"}},
			},
		}
	}

	err := getCollection(config.UserCollection).Find(query).All(&users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.APIError(err.Error()))
		return
	}

	utils.SendJSON(w, users)
}

func singleUser(w http.ResponseWriter, r *http.Request) {
	loggedInUser := middleware.GetUserSession(r)
	if loggedInUser == nil {
		loggedInUser = &models.User{}
	}

	var u models.User
	var err error

	username := mux.Vars(r)["username"]
	coll := getCollection(config.UserCollection)

	switch r.Method {
	case "GET":
		err = coll.FindId(username).One(&u)
	case "DELETE":
		if !u.IsAdmin && loggedInUser.Username != username {
			err = errors.New("not authorized")
			break
		}

		err = coll.RemoveId(username)
	case "PUT":
		if !u.IsAdmin && loggedInUser.Username != username {
			err = errors.New("not authorized")
			break
		}

		var jr map[string]models.User

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&jr)
		if err != nil {
			break
		}

		u, ok := jr["user"]
		if !ok {
			err = errors.New("invalid object schema")
			break
		}

		err = coll.Update(bson.M{"_id": username},
			bson.M{"$set": bson.M{
				"username":   u.Username,
				"email":      u.Email,
				"fullname":   u.FullName,
				"profilePic": u.ProfilePic,
			}})
	}

	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, u)
}
