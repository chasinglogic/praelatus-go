package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/urfave/cli"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func cleanDB(c *cli.Context) error {
	conn, err := mgo.Dial(config.DBURL())
	if err != nil {
		panic(err)
	}

	db := conn.DB(config.DBName())

	var r bson.M
	var ans string

	if c.Bool("yes") {
		ans = "y"
	} else {
		fmt.Println("***********************************WARNING***********************************")
		fmt.Println("This will delete ALL DATA in the database. This command is only useful for")
		fmt.Println("testing.")
		fmt.Println("********************DO NOT RUN THIS ON A PRODUCTION SYSTEM!********************")
		fmt.Print("Are you sure you want to DELETE ALL OF YOUR DATA? y/N ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ = reader.ReadString('\n')

	}

	if strings.HasPrefix(strings.ToLower(ans), "y") {
		err = db.Run("dropDatabase", &r)
		if err != nil {
			fmt.Println("ERROR:", err)
			return err
		}

		fmt.Println(r)
	}

	return nil
}

func testDB(c *cli.Context) error {
	fmt.Println("Testing connection to the database...")
	fmt.Println("URL:", config.DBURL())

	conn, err := mgo.Dial(config.DBURL())
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return err
	}

	err = conn.Ping()
	if err != nil {
		fmt.Println("Failed to ping:", err)
		return err
	}

	fmt.Println("Connection successful!")
	return nil
}

func adminCreateUser(c *cli.Context) error {
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

	conn, err := mgo.Dial(config.DBURL())
	if err != nil {
		return err
	}

	err = conn.DB(config.DBName()).C(config.UserCollection).Insert(&u)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return err
	}

	return nil
}
