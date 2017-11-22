// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/spf13/cobra"
)

var (
	isAdmin  bool
	username string
	password string
	fullName string
	email    string
)

func init() {
	createUser.Flags().StringVarP(&username, "username", "u", "",
		"Username for the new user.")
	createUser.Flags().StringVarP(&password, "password", "p", "",
		"Password for the new user.")
	createUser.Flags().StringVarP(&email, "email", "e", "",
		"Email for the new user.")
	createUser.Flags().StringVarP(&fullName, "full-name", "n", "",
		"Email for the new user.")
	createUser.Flags().BoolVarP(&isAdmin, "admin", "a", false,
		"If given the created user will be an Admin.")

	db.AddCommand(testdb)
	db.AddCommand(setupdb)
	admin.AddCommand(createUser)
}

var db = &cobra.Command{
	Use:   "db",
	Short: "Commands for interacting with the database.",
}

var testdb = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing connection to the database...")
		fmt.Println("URL:", config.DBURL())
		r := loadRepo()
		err := r.Test()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Connection failed!")
			return
		}

		fmt.Println("Connection successful!")
	},
}

var admin = &cobra.Command{
	Use:   "admin",
	Short: "Commands for administrating this instance of Praelatus.",
}

var createUser = &cobra.Command{
	Use:   "create_user",
	Short: "Create a user in the database.",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			fmt.Println("missing required --username flag")
			os.Exit(1)
		}

		if password == "" {
			fmt.Println("missing required --password flag")
			os.Exit(1)
		}

		if fullName == "" {
			fmt.Println("missing required --fullName flag")
			os.Exit(1)
		}

		if email == "" {
			fmt.Println("missing required --email flag")
			os.Exit(1)
		}

		u, err := models.NewUser(username, password, fullName, email, isAdmin)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		r := loadRepo()
		_, err = r.Users().Create(&models.User{IsAdmin: isAdmin}, *u)
		if err != nil {
			fmt.Println("Error creating user:", err)
		}
	},
}

var setupdb = &cobra.Command{
	Use:   "setup",
	Short: "Set up the database. Should only be run once on a new installation.",
	Run: func(cmd *cobra.Command, args []string) {
		r := loadRepo()
		err := r.Init()
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}
