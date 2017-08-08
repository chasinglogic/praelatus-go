package models

// Team maps directly to the teams database table.
type Team struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Lead    User   `json:"lead"`
	Members []User `json:"members,omitempty"`
}

func (t *Team) String() string {
	return jsonString(t)
}
