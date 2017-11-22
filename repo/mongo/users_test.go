// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"testing"
)

func TestUserGet(t *testing.T) {
	user, err := r.Users().Get(&admin, "testuser")
	if err != nil {
		t.Error(err)
		return
	}

	if user.Email == "" {
		t.Error("Expected an email got: ", user)
	}
}

func TestUserSearch(t *testing.T) {
	us, e := r.Users().Search(&admin, "")
	if e != nil {
		t.Error(e)
		return
	}

	if us == nil || len(us) == 0 {
		t.Error("Expected to get users instead got none.")
	}
}

func TestUserUpdate(t *testing.T) {
	u, e := r.Users().Get(&admin, "testuser")
	if e != nil {
		t.Error(e)
		return
	}

	u.Email = "saved@test.com"

	e = r.Users().Update(&admin, u.Username, u)
	if e != nil {
		t.Error(e)
		return
	}

	u2, e := r.Users().Get(&admin, "testuser")
	if e != nil {
		t.Error(e)
		return
	}

	if u2.Email != "saved@test.com" {
		t.Errorf("Expected: saved@test.com Got: %s\n", u.Email)
	}
}

func TestUserDelete(t *testing.T) {
	e := r.Users().Delete(&admin, "testuser")
	if e != nil {
		t.Error(e)
		return
	}

	if _, e = r.Users().Get(&admin, "testuser"); e == nil {
		t.Errorf("Expected an error getting user but got none.")
	}
}
