// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"log"

	"golang.org/x/crypto/bcrypt"
)

// Sanitizer is implemented by any model which needs to be sanitized before
// being JSON serialized
type Sanitizer interface {
	Sanitize() interface{}
}

// Users is an alias for a slice of Users which implements Sanitize
type Users []User

// Sanitize implements models.Sanitizer so that we don't send sensitive info
// back to the client
func (ur Users) Sanitize() interface{} {
	newUsers := make([]interface{}, len(ur))

	for i := range ur {
		newUsers[i] = ur[i].Sanitize()
	}

	return newUsers
}

// UserRole is a mapping of a user to a role
type UserRole struct {
	Role    Role   `json:"role"`
	Project string `json:"project"`
}

// User represents a user of our application
type User struct {
	Username   string   `json:"username" bson:"_id" required:"true"`
	Password   string   `json:"password,omitempty" required:"true"`
	Email      string   `json:"email" required:"true"`
	FullName   string   `json:"fullName" required:"true"`
	ProfilePic string   `json:"profilePicture" `
	IsAdmin    bool     `json:"isAdmin,omitempty"`
	IsActive   bool     `json:"isActive,omitempty"`
	Settings   Settings `json:"settings,omitempty"`

	Roles []UserRole `json:"roles"`
}

// CheckPw will verify if the given password matches for this user. Logs any
// errors it encounters
func (u User) CheckPw(pw []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), pw)
	if err == nil {
		return true
	}

	log.Println("Error checking password:", err)
	return false
}

// Sanitize implements models.Sanitizer so that we don't send sensitive info
// back to the client
func (u User) Sanitize() interface{} {
	u.Password = ""
	return u
}

func (u User) String() string {
	return jsonString(u)
}

// ProjectsMemberOf returns an array of project keys which this user has a
// role in.
func (u User) ProjectsMemberOf() []string {
	projectKeys := make([]string, len(u.Roles))

	for i := range u.Roles {
		projectKeys[i] = u.Roles[i].Project
	}

	return projectKeys
}

// RolesForProject will return an array of the roles a this user has for the
// given project.
func (u User) RolesForProject(p Project) []Role {
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
