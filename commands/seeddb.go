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
	return repo.Seed(config.LoadRepo())
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
