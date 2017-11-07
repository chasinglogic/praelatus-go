// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package event

import (
	"strings"

	"github.com/praelatus/praelatus/models"
)

// Generic is an event used when there needs to be a notification recorded but
// nothing else.
type Generic struct {
	User           models.User
	InProject      models.Project
	ActionedTicket models.Ticket
	EventType      Type
}

// ActioningUser will return the user who performed the transition
func (ge Generic) ActioningUser() models.User { return ge.User }

// Project will return the project the actioned ticket belongs to
func (ge Generic) Project() models.Project { return ge.InProject }

// Ticket will return the ticket being transitioned
func (ge Generic) Ticket() models.Ticket { return ge.ActionedTicket }

// Data returns nothing for a Generic event
func (ge Generic) Data() interface{} { return nil }

// Type will return the appropriate event type
func (ge Generic) Type() Type { return ge.EventType }

// String will return an appropriate user-readable string describing the event
func (ge Generic) String() string {
	return ge.User.Username + " " +
		strings.ToLower(string(ge.Type())) + " " + ge.ActionedTicket.Key
}
