// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"time"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type projectRepo struct {
	conn *mgo.Session
}

func (p projectRepo) coll() *mgo.Collection {
	return p.conn.DB(dbName).C(projects)
}

func permWithID(u *models.User, uid string) bson.M {
	base := permQuery(u)
	return bson.M{
		"$and": []bson.M{
			{
				"_id": uid,
			},
			base,
		},
	}
}

func (p projectRepo) Get(u *models.User, uid string) (models.Project, error) {
	var project models.Project
	q := permWithID(u, uid)
	err := p.coll().Find(q).One(&project)
	return project, mongoErr(err)
}

func (p projectRepo) Update(u *models.User, uid string, updated models.Project) error {
	q := permWithID(u, uid)
	return mongoErr(p.coll().Update(q, updated))
}

func (p projectRepo) Create(u *models.User, project models.Project) (models.Project, error) {
	if u == nil || !u.IsAdmin {
		return models.Project{}, repo.ErrAdminRequired
	}

	project.CreatedDate = time.Now()
	return project, mongoErr(p.coll().Insert(project))
}

func (p projectRepo) Delete(u *models.User, uid string) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrAdminRequired
	}

	return p.coll().RemoveId(uid)
}

func (p projectRepo) Search(u *models.User, query string) ([]models.Project, error) {
	base := permQuery(u)
	q := base

	if query != "" {
		q = bson.M{
			"$and": []bson.M{
				base,
				{
					"$or": []bson.M{
						{
							"name": bson.M{
								"$regex":   query,
								"$options": "i",
							},
						},
						{
							"key": bson.M{
								"$regex":   query,
								"$options": "i",
							},
						},
						{
							"lead": bson.M{
								"$regex":   query,
								"$options": "i",
							},
						},
					},
				},
			},
		}
	}

	var projects []models.Project
	err := p.coll().Find(q).All(&projects)
	return projects, mongoErr(err)
}

func (p projectRepo) HasLead(u *models.User, lead models.User) ([]models.Project, error) {
	base := permQuery(u)
	q := bson.M{
		"$and": []bson.M{
			base,
			{
				"lead": lead.Username,
			},
		},
	}

	var projects []models.Project
	err := p.coll().Find(q).All(&projects)
	return projects, mongoErr(err)
}
