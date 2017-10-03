// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Praelatus is an Open Source bug tracking and ticketing system. The
// backend API is written in Go and the frontend is written in Vue.js.
// You are viewing the Godoc for the Backend if you would like
// information about how to use the API as a client or how to start
// working on the backend visit http://docs.praelatus.io
// Copyright (C) 2017 Mathew Robinson <mrobinson@praelatus.io>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"os"
	"strings"

	"github.com/praelatus/praelatus/commands"
	"github.com/urfave/cli"
)

var version = "master"
var commit = ""
var date = ""

func main() {
	app := cli.NewApp()
	app.Name = "praelatus"
	app.Usage = "Praelatus, an open source bug tracker / ticketing system"

	if !strings.HasPrefix(version, "SNAPSHOT") {
		app.Version = version
	} else {
		app.Version = "DEV-" + commit
	}

	app.Action = commands.RunServer
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "devmode",
			Usage: "runs server in devmode which changes some security behavior to ease development",
		},
		cli.BoolFlag{
			Name:  "profile",
			Usage: "runs server profiling for debugging hotspots in the code",
		},
	}

	app.Authors = []cli.Author{
		{
			Name:  "Mathew Robinson",
			Email: "mrobinson@praelatus.io",
		},
		{
			Name:  "Wes Swett",
			Email: "wswett@praelatus.io",
		},
		{
			Name:  "Ryan Brzezinski",
			Email: "ryan.brzezinski867@gmail.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "seeddb",
			Usage:  "seed the database with test data",
			Action: commands.SeedDB,
		},
		{
			Name:   "serve",
			Usage:  "start running the REST api",
			Action: commands.RunServer,
		},
		{
			Name:  "config",
			Usage: "various commands for interacting with praelatus config",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "view the configuration for this instance, useful for debugging",
					Action: commands.ShowConfig,
				},
				{
					Name:   "gen",
					Usage:  "generate a config.json based on the current environment variables or defaults",
					Action: commands.GenConfig,
				},
			},
		},
		{
			Name:   "testdb",
			Usage:  "will test the connections to the databases",
			Action: commands.TestDB,
		},
		{
			Name:   "cleandb",
			Usage:  "will clean the database (remove all data), useful for testing",
			Action: commands.CleanDB,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes",
					Usage: "skip the warning prompt",
				},
			},
		},
		{
			Name:   "setupdb",
			Usage:  "Sets up indexes and other administrative tasks for the database",
			Action: commands.SetupDB,
		},
		{
			Name:  "admin",
			Usage: "various admin functions for the instance",
			Subcommands: []cli.Command{
				{
					Name:   "createUser",
					Usage:  "create a user, useful for creating admin accounts",
					Action: commands.AdminCreateUser,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "username",
						},
						cli.StringFlag{
							Name: "password",
						},
						cli.StringFlag{
							Name: "fullName",
						},
						cli.StringFlag{
							Name: "email",
						},
						cli.BoolFlag{
							Name:  "admin",
							Usage: "Indicates whether the user is a system admin",
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
