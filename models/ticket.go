package models

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TicketType represents the type of ticket.
type TicketType struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name"`
}

// Ticket represents a ticket
type Ticket struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	CreatedDate time.Time     `json:"created_date"`
	UpdatedDate time.Time     `json:"updated_date"`
	Key         string        `json:"key"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Status      string        `json:"status"`

	Type     TicketType `json:"ticket_type"`
	Reporter User       `json:"reporter"`
	Assignee User       `json:"assignee"`

	Labels []string `json:"labels"`

	Fields   []FieldValue `json:"fields"`
	Comments []Comment    `json:"comments,omitempty"`

	Workflow bson.ObjectId `json:"workflow"`
	Project  bson.ObjectId `json:"project"`
}

func (t *Ticket) String() string {
	return jsonString(t)
}

// Transition searches through the available transitions for the ticket
// returning a boolean indicating success or failure and the transition
func (t *Ticket) Transition(db *mgo.DB, name string) (Transition, bool) {
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

// Label is a label used on tickets
type Label struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name"`
}

func (l *Label) String() string {
	return jsonString(l)
}

// Comment is a comment on an issue / ticket.
type Comment struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	UpdatedDate time.Time     `json:"updated_date"`
	CreatedDate time.Time     `json:"created_date"`
	TicketKey   string        `json:"ticket_key"`
	Body        string        `json:"body"`
	Author      User          `json:"author"`
}

func (c *Comment) String() string {
	return jsonString(c)
}
