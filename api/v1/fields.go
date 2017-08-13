package v1

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
)

func fieldRouter(router *mux.Router) {
	router.HandleFunc("/fields", GetAllFields).Methods("GET")

	// router.HandleFunc("/fieldschemes", GetAllFieldSchemes)
}

// GetAllFields will retrieve all fields from the DB and send a JSON response
func GetAllFields(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		w.WriteHeader(403)
		w.Write(utils.APIError("you must be logged in to view all fields"))
		return
	}

	var fields []models.Field
	coll := Conn.DB(config.DBName()).C(config.FieldSchemeCollection)

	iter := coll.Pipe([]bson.M{
		bson.M{
			"$unwind": "$fields",
		},
		bson.M{
			"$group": bson.M{
				"$_id": bson.M{
					"name":     "$fields.name",
					"dataType": "$fields.dataType",
					"options":  "$fields.options",
				},
			},
		},
	}).Iter()

	err := iter.Err()
	if err != nil {
		w.WriteHeader(500)
		w.Write(utils.APIError(err.Error()))
		log.Println(err)
		return
	}

	var t struct {
		ID models.Field `bson:"_id,omitempty"`
	}

	for iter.Next(&t) {
		fields = append(fields, t.ID)
	}

	utils.SendJSON(w, fields)
}
