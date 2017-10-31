// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"errors"
	"strconv"
	"time"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/models/permission"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ticketRepo struct {
	conn *mgo.Session
}

func (t ticketRepo) coll() *mgo.Collection {
	return t.conn.DB(dbName).C(tickets)
}

func (t ticketRepo) Get(u *models.User, uid string) (models.Ticket, error) {
	var ticket models.Ticket
	err := t.coll().FindId(uid).One(&ticket)
	if err != nil {
		return ticket, mongoErr(err)
	}

	var p models.Project

	err = t.conn.DB(dbName).C(projects).FindId(ticket.Project).One(&p)
	if err != nil {
		return ticket, mongoErr(err)
	}

	if len(models.HasPermission(permission.ViewProject, *u, p)) == 0 {
		return ticket, repo.ErrUnauthorized
	}

	return ticket, err
}

func (t ticketRepo) Update(u *models.User, uid string, updated models.Ticket) error {
	// TODO: Have to validate field schema.
	return errors.New("unimplemented")
}

func (t ticketRepo) AddComment(u *models.User, uid string, comment models.Comment) (models.Ticket, error) {
	var ticket models.Ticket

	comment.CreatedDate = time.Now()
	comment.UpdatedDate = time.Now()

	err := t.coll().UpdateId(uid, bson.M{
		"$push": bson.M{
			"comments": comment,
		},
	})

	return ticket, mongoErr(err)
}

func (t ticketRepo) Create(u *models.User, ticket models.Ticket) (models.Ticket, error) {

	var p models.Project

	err := t.conn.DB(dbName).C(projects).FindId(ticket.Project).One(&p)
	if err != nil {
		return models.Ticket{}, mongoErr(err)
	}

	if len(models.HasPermission(permission.CreateTicket, *u, p)) == 0 {
		return models.Ticket{}, repo.ErrUnauthorized
	}

	if !p.HasTicketType(ticket.Type) {
		return models.Ticket{}, repo.ErrInvalidTicketType
	}

	var fs models.FieldScheme

	err = t.conn.DB(dbName).C(fieldSchemes).FindId(p.FieldScheme).One(&fs)
	if err != nil {
		return models.Ticket{}, mongoErr(err)
	}

	if err := fs.ValidateTicket(ticket); err != nil {
		if err.Error() == "no fields set for this ticket type and default not set" {

			return models.Ticket{}, repo.ErrInvalidFieldsForTicket
		}

		return models.Ticket{}, repo.ErrInvalidFieldsForTicket
	}

	ticket.Workflow = p.GetWorkflow(ticket.Type)

	var wkf models.Workflow

	err = t.conn.DB(dbName).C(workflows).FindId(ticket.Workflow).One(&wkf)
	if err != nil {
		return models.Ticket{}, mongoErr(err)
	}

	key, err := t.NextTicketKey(u, ticket.Project)
	if err != nil {
		return models.Ticket{}, mongoErr(err)
	}

	ticket.Key = key
	ticket.CreatedDate = time.Now()
	ticket.UpdatedDate = time.Now()
	ticket.Comments = []models.Comment{}
	ticket.Status = wkf.CreateTransition().ToStatus

	err = t.coll().Insert(ticket)
	if err != nil {
		return ticket, mongoErr(err)
	}

	return ticket, nil
}

func (t ticketRepo) Delete(u *models.User, uid string) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrUnauthorized
	}

	var ticket models.Ticket

	err := t.coll().FindId(uid).One(&ticket)
	if err != nil {
		return mongoErr(err)
	}

	var project models.Project

	err = t.conn.DB(dbName).C(projects).FindId(ticket.Project).One(&project)
	if err != nil {
		return mongoErr(err)
	}

	// TODO: Somehow move this into the query above?
	if len(models.HasPermission(permission.RemoveTicket, *u, project)) == 0 {
		return repo.ErrUnauthorized
	}

	return mongoErr(t.coll().RemoveId(uid))
}

func (t ticketRepo) Search(u *models.User, query string) ([]models.Ticket, error) {
	permQ := permQuery(u)
	var p []models.Project

	err := t.conn.DB(dbName).C(projects).Find(permQ).
		Select(bson.M{"_id": 1}).All(&p)
	if err != nil {
		return nil, mongoErr(err)
	}

	keys := make([]string, len(p))
	for i, prj := range p {
		keys[i] = prj.Key
	}

	var tickets []models.Ticket

	tQuery := bson.M{
		"project": bson.M{
			"$in": keys,
		},
	}

	if query != "" {
		tQuery = bson.M{
			"$and": []bson.M{
				tQuery,
				{
					"$or": []bson.M{
						{
							"key": bson.M{"$regex": query},
						},
						{
							"description": bson.M{"$regex": query},
						},
						{
							"summary": bson.M{"$regex": query},
						},
					},
				},
			},
		}
	}

	err = t.coll().Find(tQuery).All(&tickets)
	return tickets, err
}

func (t ticketRepo) NextTicketKey(u *models.User, projectKey string) (string, error) {
	count, err := t.coll().Find(bson.M{"project": projectKey}).Count()
	return projectKey + "-" + strconv.Itoa(count+1), mongoErr(err)
}

func (t ticketRepo) LabelSearch(u *models.User, query string) ([]string, error) {
	var labelDoc struct{ Labels []string }

	if u == nil {
		return []string{}, repo.ErrUnauthorized
	}

	pipe := t.coll().Pipe([]bson.M{
		{
			"$unwind": "$labels",
		},
		{
			"$group": bson.M{
				"_id": "$labels",
			},
		},
		{
			"$match": bson.M{
				"_id": bson.M{
					"$regex": query,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "available",
				"labels": bson.M{
					"$push": "$_id",
				},
			},
		},
	})

	err := pipe.One(&labelDoc)
	return labelDoc.Labels, mongoErr(err)
}
