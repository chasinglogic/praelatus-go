package models

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Ticket represents a ticket
type Ticket struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	CreatedDate time.Time     `json:"createdDate"`
	UpdatedDate time.Time     `json:"updatedDate"`
	Key         string        `json:"key"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Reporter    string        `json:"reporter"`
	Assignee    string        `json:"assignee"`
	Type        string        `json:"ticketType"`
	Labels      []string      `json:"labels"`

	Fields   []Field   `json:"fields"`
	Comments []Comment `json:"comments,omitempty"`

	Workflow bson.ObjectId `json:"workflow"`
	Project  bson.ObjectId `json:"project"`
}

func (t *Ticket) String() string {
	return jsonString(t)
}

// Transition searches through the available transitions for the ticket
// returning a boolean indicating success or failure and the transition
func (t *Ticket) Transition(db *mgo.Database, name string) (Transition, bool) {
	var workflow Workflow

	err := db.C("workflows").FindId(t.Workflow).One(&workflow)
	if err != nil {
		log.Println(err.Error())
		return Transition{}, false
	}

	for _, transition := range workflow.Transitions {
		if transition.Name == name && t.Status == transition.FromStatus {
			return transition, true
		}
	}

	return Transition{}, false
}

// Comment is a comment on an issue / ticket.
type Comment struct {
	UpdatedDate time.Time `json:"updatedDate"`
	CreatedDate time.Time `json:"createdDate"`
	TicketKey   string    `json:"ticketKey"`
	Body        string    `json:"body"`
	Author      string    `json:"author"`
}

func (c *Comment) String() string {
	return jsonString(c)
}
