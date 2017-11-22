// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package v1_test

import (
	"encoding/json"
	"testing"

	"github.com/praelatus/praelatus/models"
)

func ticketFromJSON(jsn []byte) (interface{}, error) {
	var tk models.Ticket
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func ticketsFromJSON(jsn []byte) (interface{}, error) {
	var tk []models.Ticket
	err := json.Unmarshal(jsn, &tk)
	return tk, err
}

func toTickets(v interface{}) []models.Ticket {
	return v.([]models.Ticket)
}

func toTicket(v interface{}) models.Ticket {
	return v.(models.Ticket)
}

var ticketRouteTests = []routeTest{
	{
		Name:         "Read Ticket",
		Endpoint:     "/api/v1/tickets/TEST-1",
		Method:       "GET",
		ExpectedCode: 200,
		Converter:    ticketFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			tk := toTicket(v)

			if tk.Key != "TEST-1" {
				t.Errorf("Expected TEST-1 Got %s", tk.Key)
			}
		},
	},

	{
		Name:         "Read All Tickets",
		Endpoint:     "/api/v1/tickets",
		Method:       "GET",
		ExpectedCode: 200,
		Converter:    ticketsFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			tk := toTickets(v)

			if len(tk) <= 1 {
				t.Errorf("Expected More Than 1 Ticket Got: %d", len(tk))
			}

			if tk[0].Key != "TEST-1" {
				t.Errorf("Expected TEST-1 Got: %s", tk[0].Key)
			}
		},
	},

	{
		Name:         "Create Ticket",
		Endpoint:     "/api/v1/tickets",
		Method:       "POST",
		Admin:        true,
		ExpectedCode: 200,
		Body: models.Ticket{
			Summary:     "A fake test ticket.",
			Description: "Not a useful description.",
			Reporter:    "testuser",
			Project:     "TEST",
			Type:        "Bug",
		},
		Converter: ticketFromJSON,
		Validator: func(v interface{}, t *testing.T) {
			tk := toTicket(v)

			if tk.Key != "" {
				t.Errorf("Expected A Ticket Key Got None\n")
			}
		},
	},

	// TODO: Implement these routes
	// {
	// 	Name:     "Add Comment",
	// 	Endpoint: "/api/v1/tickets/TEST-1/comments",
	// 	Method:   "POST",
	// 	Admin:    true,
	// 	Body: models.Comment{
	// 		Body: "THIS IS A TEST COMMENT",
	// 	},
	// 	ExpectedCode: 200,
	// 	Converter:    ticketFromJSON,
	// 	Validator: func(v interface{}, t *testing.T) {
	// 		tk := toTicket(v)

	// 		for _, comment := range tk.Comments {
	// 			if comment.Body == "THIS IS A TEST COMMENT" && comment.Author == "testadmin" {
	// 				return
	// 			}
	// 		}

	// 		t.Fail()
	// 	},
	// },

	// {
	// 	Name:      "Transition Ticket",
	// 	Login:     true,
	// 	Method:    "POST",
	// 	Endpoint:  "/api/v1/tickets/TEST-1/transition?name=In%20Progress",
	// 	Converter: ticketFromJSON,
	// 	Validator: func(v interface{}, t *testing.T) {
	// 		tk := toTicket(v)

	// 		if tk.Status != "In Progress" {
	// 			t.Errorf("Expected In Progress Got: %s", tk.Status)
	// 		}
	// 	},
	// },
}

func TestTicketRoutes(t *testing.T) {
	testRoutes(ticketRouteTests, t)
}
