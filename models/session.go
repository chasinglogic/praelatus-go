package models

import "time"

// Session stores a user with the expiration time of the session
type Session struct {
	ID      string `bson:"_id"`
	Expires time.Time
	User    User
}

func (s Session) String() string {
	return jsonString(s)
}
