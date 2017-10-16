// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import "gopkg.in/mgo.v2/bson"

// Workflow is the container for issues and keeps track of available transitions
type Workflow struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" required:"true"`
	Transitions []Transition  `json:"transitions" required:"true"`
}

func (w Workflow) String() string {
	return jsonString(w)
}

// CreateTransition will return the transition to perform on a ticket during creation
func (w Workflow) CreateTransition() Transition {
	for _, t := range w.Transitions {
		if t.FromStatus == "Create" {
			return t
		}
	}

	return Transition{ToStatus: "null"}
}

// Transition contains information about what hooks to perform when performing
// a transition
type Transition struct {
	Name       string `json:"name"`
	FromStatus string `json:"fromStatus"`
	ToStatus   string `json:"toStatus"`
	Hooks      []Hook `json:"hooks"`
}

func (t Transition) String() string {
	return jsonString(t)
}

// Hook contains information about what webhooks to fire when a given
// transition is run.
type Hook struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Body     string `json:"body"`
}

func (h Hook) String() string {
	return jsonString(h)
}
