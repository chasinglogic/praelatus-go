// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	"github.com/praelatus/praelatus/repo/mongo"
	"gopkg.in/mgo.v2/bson"
)

var r repo.Repo
var c repo.Cache

var admin = models.User{
	Username: "testadmin",
	IsAdmin:  true,
}

var user = models.User{
	Username: "testuser",
	IsAdmin:  false,
}

var fsID bson.ObjectId
var wID bson.ObjectId

func init() {
	r = mongo.New("mongodb://localhost/praelatus_test")
	e := r.Clean()
	if e != nil {
		panic(e)
	}

	e = repo.Seed(r)
	if e != nil {
		panic(e)
	}

	c = mongo.NewCache("mongodb://localhost/praelatus_test")

	mr, ok := r.(mongo.Repo)
	if !ok {
		panic("Should be a mongo repo.")
	}

	var f models.FieldScheme

	e = mr.Conn.DB("praelatus").C("field_schemes").Find(nil).One(&f)
	if e != nil {
		panic(e)
	}

	fsID = f.ID

	var w models.Workflow

	e = mr.Conn.DB("praelatus").C("workflows").Find(nil).One(&w)
	if e != nil {
		panic(e)
	}

	wID = w.ID
}
