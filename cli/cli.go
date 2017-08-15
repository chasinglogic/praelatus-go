// Package cli contains all of the functions used in our cli, it additionally
// is where the cli itsself is defined and ran
package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// Run runs the cli of Praelatus with the given argv
func Run(args []string) {
	app := cli.NewApp()
	app.Name = "praelatus"
	app.Usage = "Praelatus, an Open Source bug tracker / ticketing system"
	app.Version = "0.2.0"
	app.Action = runServer
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
			Action: seedDB,
		},
		{
			Name:   "serve",
			Usage:  "start running the REST api",
			Action: runServer,
		},
		{
			Name:  "config",
			Usage: "various commands for interacting with praelatus config",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "view the configuration for this instance, useful for debugging",
					Action: showConfig,
				},
				{
					Name:   "gen",
					Usage:  "generate a config.json based on the current environment variables or defaults",
					Action: genConfig,
				},
			},
		},
		{
			Name:   "testdb",
			Usage:  "will test the connections to the databases",
			Action: testDB,
		},
		{
			Name:   "cleandb",
			Usage:  "will clean the database (remove all data), useful for testing",
			Action: cleanDB,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes",
					Usage: "skip the warning prompt",
				},
			},
		},
		{
			Name:   "setupdb",
			Usage:  "Sets up indexes and other such items on the database.",
			Action: setupDB,
		},
		{
			Name:  "admin",
			Usage: "various admin functions for the instance",
			Subcommands: []cli.Command{
				{
					Name:   "createUser",
					Usage:  "create a user, useful for creating admin accounts",
					Action: adminCreateUser,
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
							Usage: "when this flag is given user will be created as an system admin",
						},
					},
				},
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
	}
}
