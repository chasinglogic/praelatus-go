// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package event

import "github.com/praelatus/praelatus/models"

// Transition should be fired whenever a ticket is transitioned through a
// workflow
type Transition struct {
	User           models.User
	InProject      models.Project
	ActionedTicket models.Ticket
	Transition     models.Transition
}

// ActioningUser will return the user who performed the transition
func (te Transition) ActioningUser() models.User { return te.User }

// Project will return the project the actioned ticket belongs to
func (te Transition) Project() models.Project { return te.InProject }

// Ticket will return the ticket being transitioned
func (te Transition) Ticket() models.Ticket { return te.ActionedTicket }

// Data will return the transition being performed
func (te Transition) Data() interface{} { return te.Transition }

// Type will return the appropriate event type
func (te Transition) Type() Type { return TransitionEvent }

// String will return an appropriate user-readable string describing the event
func (te Transition) String() string {
	s := te.User.Username + " transitioned " + te.ActionedTicket.Key +
		" to " + te.Transition.ToStatus.Name

	if te.Transition.FromStatus.Type != models.StatusNull {
		s += " from " + te.Transition.FromStatus.Name
	}

	return s
}
