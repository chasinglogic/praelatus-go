package cli

// import (
// 	"github.com/praelatus/backend/config"
// 	"github.com/praelatus/backend/models"
// 	"github.com/urfave/cli"
// )

// func adminCreateUser(c *cli.Context) error {
// 	username := c.String("username")
// 	password := c.String("password")
// 	fullName := c.String("fullName")
// 	email := c.String("email")
// 	admin := c.Bool("admin")

// 	if username == "" {
// 		return cli.NewExitError("missing required --username flag", 1)
// 	}

// 	if password == "" {
// 		return cli.NewExitError("missing required --password flag", 1)
// 	}

// 	if fullName == "" {
// 		return cli.NewExitError("missing required --fullName flag", 1)
// 	}

// 	if email == "" {
// 		return cli.NewExitError("missing required --email flag", 1)
// 	}

// 	u, err := models.NewUser(username, password, fullName, email, admin)
// 	if err != nil {
// 		return cli.NewExitError(err.Error(), 1)
// 	}

// 	s := config.Store()
// 	if err := s.Users().New(u); err != nil {
// 		return cli.NewExitError(err.Error(), 1)
// 	}

// 	return nil
// }
