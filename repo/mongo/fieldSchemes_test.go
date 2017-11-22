// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import "testing"

func TestFieldSchemeGet(t *testing.T) {
	t.Log(fsID)

	fieldScheme, err := r.Fields().Get(&admin, fsID.Hex())
	if err != nil {
		t.Error(err)
		return
	}

	if fieldScheme.Name == "" {
		t.Error("Expected a name got: ", fieldScheme.Name)
	}
}

func TestFieldSchemeSearch(t *testing.T) {
	fs, e := r.Fields().Search(&admin, "")
	if e != nil {
		t.Error(e)
		return
	}

	if fs == nil || len(fs) == 0 {
		t.Error("Expected to get field schemes instead got none.")
	}
}

func TestFieldSchemeUpdate(t *testing.T) {
	f, e := r.Fields().Get(&admin, fsID.Hex())
	if e != nil {
		t.Error(e)
		return
	}

	f.Name = "Test field scheme save"

	e = r.Fields().Update(&admin, fsID.Hex(), f)
	if e != nil {
		t.Error(e)
		return
	}

	f2, e := r.Fields().Get(&admin, fsID.Hex())
	if e != nil {
		t.Error(e)
		return
	}

	if f2.Name != "Test field scheme save" {
		t.Errorf("Expected: Test field scheme save Got: %s\n", f.Name)
	}
}

func TestFieldSchemeDelete(t *testing.T) {
	e := r.Fields().Delete(&admin, fsID.Hex())
	if e != nil {
		t.Error(e)
		return
	}

	if _, e = r.Fields().Get(&admin, fsID.Hex()); e == nil {
		t.Errorf("Expected an error getting field scheme but got none.")
	}
}
