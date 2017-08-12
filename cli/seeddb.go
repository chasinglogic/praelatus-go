package cli

import (
	"time"

	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
	"github.com/urfave/cli"
	mgo "gopkg.in/mgo.v2"
)

func seedDB(c *cli.Context) error {
	conn, err := mgo.Dial(config.DBURL())
	if err != nil {
		panic(err)
	}
	db := conn.DB(config.DBName())

	users := []models.User{
		models.NewUser("testadmin", "test", "Test Testerson", "test@example.com", true),
		models.NewUser("testuser", "test", "Test Testerson II", "test@example.com", false),
	}

	for _, u := range users {
		db.C("users").Insert(&u)
	}

	p := models.Project{
		Key:        "TEST",
		Name:       "Test Project",
		CreateDate: time.Now(),
		Lead:       "testadmin",
		TicketTypes: []string{
			"Epic",
			"Story",
			"Bug",
			"Feature Request",
		},

		FieldScheme: map[string][]Field{
			"Story": []Field{
				{
					Name:     "Story Points",
					DataType: "INT",
				},
			},
			"": []Field{},
		},

		Permissions: map[Role][]Permission{
			"Administrator": []Permission{
				"VIEW_PROJECT",
				"ADMIN_PROJECT",
				"CREATE_TICKET",
				"COMMENT_TICKET",
				"REMOVE_COMMENT",
				"REMOVE_OWN_COMMENT",
				"EDIT_OWN_COMMENT",
				"EDIT_COMMENT",
				"TRANSITION_TICKET",
				"EDIT_TICKET",
				"REMOVE_TICKET",
			},
			"Contributor": []Permission{
				"VIEW_PROJECT",
				"CREATE_TICKET",
				"COMMENT_TICKET",
				"REMOVE_OWN_COMMENT",
				"EDIT_OWN_COMMENT",
				"TRANSITION_TICKET",
				"EDIT_TICKET",
			},
			"User": []Permission{
				"VIEW_PROJECT",
				"CREATE_TICKET",
				"COMMENT_TICKET",
			},
			"Anonymous": []Permission{
				"VIEW_PROJECT",
			},
		},
	}

	for _, dataType := range models.DataTypes {
		f := models.Field{
			Name:     "Test " + dataType + " Field",
			DataType: dataType,
		}

		if dataType == "OPT" {
			f.Options = []string{
				"High",
				"Medium",
				"Low",
			}
		}

		p.FieldScheme[""] = append(p.FieldScheme[""], f)
	}

	db.C("projects").insert(&p)

	return nil
}
