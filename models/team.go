package models

import "gopkg.in/mgo.v2/bson"

// Team maps directly to the teams database table.
type Team struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"name"`
	Lead    User          `json:"lead"`
	Members []User        `json:"members,omitempty"`
}

func (t *Team) String() string {
	return jsonString(t)
}
