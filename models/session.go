package models

// Session ties a token to a specific client
type Session struct {
	Token    string `bson:"_id"`
	ClientID string
}

func (s Session) String() string {
	return jsonString(s)
}
