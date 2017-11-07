// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"fmt"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type notificationRepo struct {
	conn *mgo.Session
}

func (nr notificationRepo) coll() *mgo.Collection {
	return nr.conn.DB(dbName).C(notifications)
}

func (nr notificationRepo) Create(u *models.User, n models.Notification) (models.Notification, error) {
	n.ID = bson.NewObjectId()
	err := nr.coll().Insert(n)
	return n, mongoErr(err)
}

func (nr notificationRepo) MarkRead(u *models.User, uid string) error {
	return mongoErr(nr.coll().Update(
		bson.M{
			"_id":     uid,
			"watcher": u.Username,
		},
		bson.M{
			"$set": bson.M{"read": true},
		}))
}

func (nr notificationRepo) ForProject(u *models.User, project models.Project, onlyUnread bool, last int) ([]models.Notification, error) {
	keys, err := getKeysUserHasPermissionTo(u, nr.conn)
	if err != nil {
		return nil, err
	}

	conditions := []bson.M{
		{
			"project": bson.M{"$in": keys},
		},
		{
			"project": project.Key,
		},
	}

	if onlyUnread {
		conditions = append(conditions, bson.M{"unread": true})
	}

	q := bson.M{
		"$and": conditions,
	}

	query := nr.coll().Find(q).Sort("-createddate")

	if last > 0 {
		query.Limit(last)
	}

	var notifications []models.Notification
	err = query.All(&notifications)

	fmt.Println(notifications, err)
	return notifications, mongoErr(err)
}

func (nr notificationRepo) ForUser(u *models.User, user models.User, onlyUnread bool, last int) ([]models.Notification, error) {
	if u.Username != user.Username && !u.IsAdmin {
		return nil, repo.ErrUnauthorized
	}

	keys, err := getKeysUserHasPermissionTo(u, nr.conn)
	if err != nil {
		return nil, err
	}

	conditions := []bson.M{
		{
			"project": bson.M{"$in": keys},
		},
		{
			"user": user.Username,
		},
	}

	if onlyUnread {
		conditions = append(conditions, bson.M{"unread": true})
	}

	q := bson.M{
		"$and": conditions,
	}

	query := nr.coll().Find(q).Sort("-createddate")

	if last > 0 {
		query.Limit(last)
	}

	var notifications []models.Notification
	err = query.All(&notifications)

	return notifications, mongoErr(err)
}
