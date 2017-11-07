// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package event

import "github.com/praelatus/praelatus/models"

// Comment should be fired whenever a ticket is transitioned through a
// workflow
type Comment struct {
	User           models.User
	InProject      models.Project
	ActionedTicket models.Ticket
	Comment        models.Comment
}

// ActioningUser will return the user who performed the transition
func (ce Comment) ActioningUser() models.User { return ce.User }

// Project will return the project the actioned ticket belongs to
func (ce Comment) Project() models.Project { return ce.InProject }

// Ticket will return the ticket being transitioned
func (ce Comment) Ticket() models.Ticket { return ce.ActionedTicket }

// Data will return the transition being performed
func (ce Comment) Data() interface{} { return ce.Comment }

// Type will return the appropriate event type
func (ce Comment) Type() Type { return CommentEvent }

// String will return an appropriate user-readable string describing the event
func (ce Comment) String() string {
	return ce.User.Username + " commented on " + ce.ActionedTicket.Key
}
