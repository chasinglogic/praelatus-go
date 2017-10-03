// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/urfave/cli"
)

// TestDB will test the connection to the configured database.
func TestDB(c *cli.Context) error {
	fmt.Println("Testing connection to the database...")
	fmt.Println("URL:", config.DBURL())
	r := config.LoadRepo()
	err := r.Test()
	if err != nil {
		fmt.Println("Connection failed!", err)
		return err
	}

	fmt.Println("Connection successful!")
	return nil
}

// AdminCreateUser will allow an admin to create a user from the server.
func AdminCreateUser(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")
	fullName := c.String("fullName")
	email := c.String("email")
	admin := c.Bool("admin")

	if username == "" {
		return cli.NewExitError("missing required --username flag", 1)
	}

	if password == "" {
		return cli.NewExitError("missing required --password flag", 1)
	}

	if fullName == "" {
		return cli.NewExitError("missing required --fullName flag", 1)
	}

	if email == "" {
		return cli.NewExitError("missing required --email flag", 1)
	}

	u, err := models.NewUser(username, password, fullName, email, admin)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	r := config.LoadRepo()
	_, err = r.Users().Create(&models.User{IsAdmin: true}, *u)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return err
	}

	return nil
}
