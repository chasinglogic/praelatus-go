// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/repo"
	"github.com/urfave/cli"
)

// SeedDB will seed the database with test data.
func SeedDB(c *cli.Context) error {
	fmt.Println("Connecting to database...")
	r := config.LoadRepo()

	fmt.Println("Seeding database with test data...")
	err := repo.Seed(r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done!")
	return err
}

// CleanDB will remove all data from the database. Useful for testing.
func CleanDB(c *cli.Context) error {
	var ans string

	if c.Bool("yes") {
		ans = "y"
	} else {
		fmt.Println("***********************************WARNING***********************************")
		fmt.Println(`This will delete ALL DATA in the database. This command is only useful for")
testing.`)
		fmt.Println("********************DO NOT RUN THIS ON A PRODUCTION SYSTEM!********************")
		fmt.Print("Are you sure you want to DELETE ALL OF YOUR DATA? y/N ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ = reader.ReadString('\n')

	}

	r := config.LoadRepo()
	if strings.HasPrefix(strings.ToLower(ans), "y") {
		return r.Clean()
	}

	return nil
}
