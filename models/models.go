// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package models contains all of our models and utility functions for
// creating and modifying them
package models

import (
	"encoding/json"
	"fmt"
)

func jsonString(i interface{}) string {
	b, e := json.MarshalIndent(i, "", "\t")
	if e != nil {
		fmt.Println(e)
		return ""
	}

	return string(b)
}
