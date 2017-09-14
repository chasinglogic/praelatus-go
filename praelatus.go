// Praelatus is an Open Source bug tracking and ticketing system. The
// backend API is written in Go and the frontend is a React.js
// app. You are viewing the Godoc for the API if you would like
// information about how to use the API as a client or how to start
// working on the backend visit https://docs.praelatus.io
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
