// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"testing"

	"github.com/praelatus/praelatus/models"
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
		Token:    "faketoken",
		ClientID: "fake",
	}

	err := c.SetSession(sess.Token, sess)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	d, err := c.GetSession(sess.Token)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	if sess.Token != d.Token || sess.ClientID != d.ClientID {
		t.Errorf("Expected d and sess to match Got %s\n", d)
		t.Log("Token match", sess.Token == d.Token)
		t.Log("ClientID match", sess.ClientID != d.ClientID)
	}

	err = c.RemoveSession(sess.Token)
	if err != nil {
		t.Errorf("Failed With Error: %s\n", err.Error())
		return
	}

	_, err = c.GetSession(sess.Token)
	if err == nil {
		t.Errorf("Expected no session to return, instead failed to remove.")
	}
}
