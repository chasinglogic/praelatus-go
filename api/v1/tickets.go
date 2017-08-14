package v1

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/api/utils"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
	"github.com/praelatus/backend/models/permission"
	"gopkg.in/mgo.v2/bson"
)

func ticketRouter(router *mux.Router) {
	router.HandleFunc("/tickets", getAllTickets).Methods("GET")
	// router.HandleFunc("/tickets", createTicket).Methods("POST")
	router.HandleFunc("/tickets/{key}", singleTicket)

	// router.HandleFunc("/tickets/{key}/addComment", addComment).Methods("POST")
}

func singleTicket(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var t models.Ticket
	var err error

	id := mux.Vars(r)["key"]
	coll := getCollection(config.TicketCollection)

	switch r.Method {
	case "GET":
		err = coll.FindId(id).One(&t)
	case "DELETE":
		err = coll.RemoveId(id)
	case "PUT":
		// TODO: Have to validate field schema.
		break
	}

	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, t)
}

// getAllTickets will return all tickets which the user has permissions to.
// TODO: Make this faster checking permissions requires multiple nested loops
// and memory allocations and so is the slowest call at 10ms which isn't
// horrible but certainly leaves much to be desired.
func getAllTickets(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	query := bson.M{
		"$or": []bson.M{
			{
				"key": bson.M{
					"$in": u.ProjectsMemberOf(),
				},
			},
			{
				"permissions.Anonymous": permission.ViewProject,
			},
		},
	}

	var projects []models.Project

	err := getCollection(config.ProjectCollection).Find(query).
		Select(bson.M{"permissions": 1, "key": 1}).All(&projects)
	if err != nil {
		log.Println("Error:", err.Error())
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	projects = models.HasPermission(permission.ViewProject, *u, projects...)

	keys := make([]string, len(projects))

	for i := range projects {
		keys[i] = projects[i].Key
	}

	var tickets []models.Ticket

	tQuery := bson.M{
		"project": bson.M{
			"$in": keys,
		},
	}

	err = getCollection(config.TicketCollection).Find(tQuery).All(&tickets)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSONR(w, models.JSONRepr{"tickets": tickets})
}
