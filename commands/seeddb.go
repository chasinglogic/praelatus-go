// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// +build !release

package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/repo"
	"github.com/spf13/cobra"
)

var skipPrompt bool

func init() {
	cleandb.Flags().BoolVarP(&skipPrompt, "yes", "y", false,
		"Skip the warning prompt when cleaning database.")

	db.AddCommand(seeddb)
	db.AddCommand(cleandb)
}

var seeddb = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with test data.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connecting to database...")
		r := config.LoadRepo()

		fmt.Println("Seeding database with test data...")
		err := repo.Seed(r)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Done!")
	},
}

var cleandb = &cobra.Command{
	Use:   "clean",
	Short: "Remove all data from the database. DO NOT RUN IN PRODUCTION.",
	Run: func(cmd *cobra.Command, args []string) {
		var ans string

		if skipPrompt {
			ans = "y"
		} else {
			fmt.Println(`
***********************************WARNING************************************
This will delete ALL DATA in the database. This command is only useful for"
testing.
********************DO NOT RUN THIS ON A PRODUCTION SYSTEM********************`)
			fmt.Print("Are you sure you want to DELETE ALL OF YOUR DATA? y/N ")
			reader := bufio.NewReader(os.Stdin)
			ans, _ = reader.ReadString('\n')

		}

		r := config.LoadRepo()
		if strings.HasPrefix(strings.ToLower(ans), "y") {
			err := r.Clean()
			if err != nil {
				fmt.Println("ERROR:", err)
			}
		}

	},
}
