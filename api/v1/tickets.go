package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

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
	router.HandleFunc("/tickets", createTicket).Methods("POST")
	router.HandleFunc("/tickets/{key}", singleTicket)

	router.HandleFunc("/tickets/{key}/addComment", addComment).Methods("POST")
}

func createTicket(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		u = &models.User{}
	}

	var tjson map[string]models.Ticket

	var t models.Ticket
	var p models.Project

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tjson)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	t = tjson["ticket"]

	projects := getCollection(config.ProjectCollection)
	err = projects.FindId(t.Project).One(&p)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(models.HasPermission(permission.CreateTicket, *u, p)) == 0 {
		utils.APIErr(w, http.StatusUnauthorized, "you do not have the required permission")
		return
	}

	if !p.HasTicketType(t.Type) {
		utils.APIErr(w, http.StatusForbidden, "not a valid ticket type for project "+t.Project)
		return
	}

	var fs models.FieldScheme

	fieldSchemes := getCollection(config.FieldSchemeCollection)
	err = fieldSchemes.FindId(p.FieldScheme).One(&fs)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := fs.ValidateTicket(t); err != nil {
		utils.APIErr(w, http.StatusBadRequest, err.Error())
		return
	}

	var wkf models.Workflow

	workflows := getCollection(config.WorkflowCollection)
	err = workflows.FindId(t.Workflow).One(&wkf)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	tickets := getCollection(config.TicketCollection)
	count, err := tickets.Find(bson.M{"project": t.Project}).Count()
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	t.Key = t.Project + "-" + strconv.Itoa(count+1)
	t.CreatedDate = time.Now()
	t.UpdatedDate = time.Now()
	t.Comments = []models.Comment{}
	t.Status = wkf.CreateTransition().ToStatus
	t.Workflow = p.GetWorkflow(t.Type)

	err = tickets.Insert(t)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, t)
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
		if err != nil {
			break
		}

		var p models.Project
		projects := getCollection(config.ProjectCollection)
		err = projects.FindId(t.Project).One(&p)
		if err != nil {
			break
		}

		if len(models.HasPermission(permission.ViewProject, *u, p)) == 0 {
			err = errors.New("unauthorized")
		}
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

	tQuery := bson.M{
		"project": bson.M{
			"$in": keys,
		},
	}

	var tickets []models.Ticket

	err = getCollection(config.TicketCollection).Find(tQuery).All(&tickets)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, tickets)
}

func addComment(w http.ResponseWriter, r *http.Request) {
	u := middleware.GetUserSession(r)
	if u == nil {
		utils.APIErr(w, http.StatusForbidden, "you must be logged in to comment")
		return
	}

	var cjson map[string]models.Comment

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cjson)
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	key := mux.Vars(r)["key"]
	tickets := getCollection(config.TicketCollection)

	c := cjson["comment"]
	c.CreatedDate = time.Now()
	c.UpdatedDate = time.Now()

	err = tickets.UpdateId(key, bson.M{
		"$push": bson.M{
			"comments": c,
		},
	})
	if err != nil {
		utils.APIErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSON(w, map[string]string{})
}
