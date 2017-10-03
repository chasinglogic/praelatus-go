// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"testing"
	"time"

	"github.com/praelatus/praelatus/models"
	"gopkg.in/mgo.v2/bson"
)

func TestSetRemoveAndGet(t *testing.T) {
	data := "hello world"

	err := c.Set("testdata", data)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	d, err := c.Get("testdata")
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	td, ok := d.(string)
	if !ok {
		t.Errorf("Expected map[string]string Got: %T\n", d)
		return
	}

	if td != "hello world" {
		t.Errorf("Expected hello world Got %s\n", td)
		t.Log(td)
		return
	}

	err = c.Remove("testdata")
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	_, err = c.Get("testdata")
	if err == nil {
		t.Errorf("Expected an error got data back.")
		return
	}
}

func TestSetRemoveAndGetSession(t *testing.T) {
	sess := models.Session{
		ID:      string(bson.NewObjectId()),
		Expires: time.Now(),
		User:    models.User{Username: "testuser"},
	}

	err := c.SetSession(sess.ID, sess)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	d, err := c.GetSession(sess.ID)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	if sess.ID != d.ID || sess.Expires.Sub(d.Expires) == 0 || sess.User.Username != d.User.Username {
		t.Errorf("Expected d and sess to match Got %s\n", d)
		t.Log("ID match", sess.ID == d.ID)
		t.Log("Expires match", sess.Expires.Sub(d.Expires) == 0)
		t.Log("Username match", sess.User.Username == d.User.Username)
	}
}
