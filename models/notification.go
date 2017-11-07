// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Notification is any kind of activity on a project for which a user might want
// an auditble history of
type Notification struct {
	ID             bson.ObjectId `bson:"_id" json:"id"`
	Project        string        `json:"project"`
	ActioningUser  string        `json:"actioningUser"`
	ActionedTicket string        `json:"actionedTicket"`
	Type           string        `json:"eventType"`
	Body           string        `json:"body"`
	Read           bool          `json:"read"`
	Watcher        string        `json:"watcher"`
	CreatedDate    time.Time     `json:"createdDate"`
}
