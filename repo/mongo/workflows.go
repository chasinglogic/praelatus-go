// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type workflowRepo struct {
	conn *mgo.Session
}

func (w workflowRepo) coll() *mgo.Collection {
	return w.conn.DB(dbName).C(workflows)
}

func (w workflowRepo) Get(u *models.User, uid string) (models.Workflow, error) {
	if u == nil {
		return models.Workflow{}, repo.ErrLoginRequired
	}

	var wkf models.Workflow
	err := w.coll().FindId(bson.ObjectIdHex(uid)).One(&wkf)

	return wkf, mongoErr(err)
}

func (w workflowRepo) Update(u *models.User, uid string, updated models.Workflow) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrAdminRequired
	}

	// FIXME: Handle what to do with tickets and projects associated with this
	// workflow

	return mongoErr(w.coll().UpdateId(bson.ObjectIdHex(uid), updated))
}

func (w workflowRepo) Create(u *models.User, workflow models.Workflow) (models.Workflow, error) {
	if u == nil || !u.IsAdmin {
		return models.Workflow{}, repo.ErrAdminRequired
	}

	workflow.ID = bson.NewObjectId()

	err := w.coll().Insert(workflow)
	return workflow, mongoErr(err)
}

func (w workflowRepo) Delete(u *models.User, uid string) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrAdminRequired
	}

	// FIXME: Handle what to do with tickets and projects associated with this
	// workflow

	return mongoErr(w.coll().RemoveId(bson.ObjectIdHex(uid)))
}

func (w workflowRepo) Search(u *models.User, query string) ([]models.Workflow, error) {
	if u == nil || !u.IsAdmin {
		return []models.Workflow{}, repo.ErrAdminRequired
	}

	var ws []models.Workflow

	q := bson.M{}
	if query != "" {
		q = bson.M{"name": bson.M{"$regex": query, "$options": "i"}}
	}

	err := w.coll().Find(q).All(&ws)
	return ws, mongoErr(err)
}
