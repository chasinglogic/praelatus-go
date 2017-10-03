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

type userRepo struct {
	conn *mgo.Session
}

func (ur userRepo) coll() *mgo.Collection {
	return ur.conn.DB(dbName).C(users)
}

func (ur userRepo) Get(u *models.User, uid string) (models.User, error) {
	var user models.User
	err := ur.coll().FindId(uid).One(&user)
	return user, mongoErr(err)
}

func (ur userRepo) Update(u *models.User, uid string, updated models.User) error {
	if u == nil || (!u.IsAdmin && u.Username != uid) {
		return repo.ErrUnauthorized
	}

	return mongoErr(ur.coll().UpdateId(uid, bson.M{
		"$set": bson.M{
			"username":   updated.Username,
			"email":      updated.Email,
			"profilepic": updated.ProfilePic,
			"fullname":   updated.FullName,
			"settings":   updated.Settings,
		},
	}))
}

func (ur userRepo) Create(u *models.User, user models.User) (models.User, error) {
	if u == nil || !u.IsAdmin {
		user.IsAdmin = false
	}

	return user, mongoErr(ur.coll().Insert(user))
}

func (ur userRepo) Delete(u *models.User, uid string) error {
	if u == nil || (!u.IsAdmin && u.Username != uid) {
		return repo.ErrUnauthorized
	}

	return mongoErr(ur.coll().RemoveId(uid))
}

func (ur userRepo) Search(u *models.User, query string) ([]models.User, error) {
	var users []models.User

	q := bson.M{}
	if query != "" {
		q = bson.M{
			"$or": []bson.M{
				{"username": bson.M{"$regex": q, "$options": "i"}},
				{"email": bson.M{"$regex": q, "$options": "i"}},
				{"fullname": bson.M{"$regex": q, "$options": "i"}},
			},
		}
	}

	err := ur.coll().Find(q).All(&users)
	return users, mongoErr(err)
}
