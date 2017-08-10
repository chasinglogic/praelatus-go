package models

import (
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
	Fields      []FieldValue  `json:"fields"`
	Labels      []Label       `json:"labels"`
	Type        TicketType    `json:"ticket_type"`
	Reporter    User          `json:"reporter"`
	Assignee    User          `json:"assignee"`
	Status      Status        `json:"status"`

	WorkflowID int64 `json:"workflow_id"`

	Transitions []Transition `json:"transitions"`
	Comments    []Comment    `json:"comments,omitempty"`

	Project Project `json:"project"`
}

func (t *Ticket) String() string {
	return jsonString(t)
}

// Transition searches through the available transitions for the ticket
// returning a boolean indicating success or failure and the transition
func (t *Ticket) Transition(name string) (Transition, bool) {
	for _, transition := range t.Transitions {
		if transition.Name == name {
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
