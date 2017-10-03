// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

// Session ties a token to a specific client
type Session struct {
	Token    string `bson:"_id"`
	ClientID string
}

func (s Session) String() string {
	return jsonString(s)
}
