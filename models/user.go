package models

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserRole struct {
	Role    Role   `json:"role"`
	Project string `json:"project"`
}

// User represents a user of our application
type User struct {
	Username   string   `json:"username" bson:"_id"`
	Password   string   `json:"password,omitempty"`
	Email      string   `json:"email"`
	FullName   string   `json:"fullName"`
	ProfilePic string   `json:"profilePicture"`
	IsAdmin    bool     `json:"isAdmin,omitempty"`
	IsActive   bool     `json:"isActive,omitempty"`
	Settings   Settings `json:"settings,omitempty"`

	Roles []UserRole
}

// CheckPw will verify if the given password matches for this user. Logs any
// errors it encounters
func (u *User) CheckPw(pw []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), pw)
	if err == nil {
		return true
	}

	log.Println("Error checking password:", err)
	return false
}

func (u *User) String() string {
	return jsonString(u)
}

// ProjectsMemberOf returns an array of project keys which this user has a
// role in.
func (u *User) ProjectsMemberOf() []string {
	projectKeys := make([]string, len(u.Roles))

	i := 0
	for _, r := range u.Roles {
		projectKeys[i] = r.Project
		i++
	}

	return projectKeys
}

func (u *User) RolesForProject(p Project) []Role {
	roles := make([]Role, 0)

	for _, r := range u.Roles {
		if r.Project == p.Key {
			roles = append(roles, r.Role)
		}
	}

	return roles
}

// NewUser will create the user after encrypting the password with bcrypt
func NewUser(username, password, fullName, email string, admin bool) (*User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}

	emailHash := md5.Sum([]byte(strings.ToLower(email)))
	eh := hex.EncodeToString(emailHash[:16])

	return &User{
		Username:   username,
		Password:   string(pw),
		Email:      email,
		FullName:   fullName,
		ProfilePic: "https://www.gravatar.com/avatar/" + eh,
		IsAdmin:    admin,
	}, nil
}

// Settings represents an individual users preferences
type Settings struct {
	DefaultProject string `json:"defaultProject,omitempty"`
	DefaultView    string `json:"defaultView,omitempty"`
}
