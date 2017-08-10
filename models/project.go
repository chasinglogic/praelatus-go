package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Project is the model used to represent a project in the database.
type Project struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	CreatedDate time.Time     `json:"created_date"`
	Name        string        `json:"name"`
	Key         string        `json:"key"`
	Lead        User          `json:"lead"`
	Homepage    string        `json:"homepage,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	Repo        string        `json:"repo,omitempty"`
}

func (p *Project) String() string {
	return jsonString(p)
}
