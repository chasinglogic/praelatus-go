package repo

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/models/permission"
)

// Seed will fill the given repo with test data.
func Seed(r Repo) error {
	var err error

	u1, _ := models.NewUser("testadmin", "test", "Test Testerson", "test@example.com", true)
	u2, _ := models.NewUser("testuser", "test", "Test Testerson II", "test@example.com", false)
	users := []models.User{
		*u1,
		*u2,
	}

	for _, u := range users {
		_, err = r.Users().Create(&models.User{IsAdmin: true}, u)
		if err != nil {
			return err
		}
	}

	fs := models.FieldScheme{
		Name: "Test Field Scheme",
		Fields: map[string][]models.Field{
			"Story": []models.Field{
				{
					Name:     "Story Points",
					DataType: "INT",
				},
			},
			"": []models.Field{},
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

		fs.Fields[""] = append(fs.Fields[""], f)
	}

	fs, err = r.Fields().Create(u1, fs)
	if err != nil {
		return err
	}

	workflows := []models.Workflow{
		{
			Name: "Test Simple Workflow",
			Transitions: []models.Transition{
				{
					Name:       "In Progress",
					FromStatus: "",
					ToStatus:   "In Progress",
					Hooks:      []models.Hook{},
				},
				{
					Name:       "Done",
					FromStatus: "",
					ToStatus:   "Done",
					Hooks:      []models.Hook{},
				},
				{
					Name:       "Backlog",
					FromStatus: "Create",
					ToStatus:   "Backlog",
					Hooks:      []models.Hook{},
				},
			},
		},
		{
			Name: "Test One Way Workflow",
			Transitions: []models.Transition{
				{
					Name:       "In Progress",
					FromStatus: "Backlog",
					ToStatus:   "In Progress",
					Hooks:      []models.Hook{},
				},
				{
					Name:       "Done",
					FromStatus: "In Progress",
					ToStatus:   "Done",
					Hooks:      []models.Hook{},
				},
				{
					Name:       "Backlog",
					FromStatus: "Create",
					ToStatus:   "Backlog",
					Hooks:      []models.Hook{},
				},
			},
		},
	}

	for i := range workflows {
		workflows[i], err = r.Workflows().Create(u1, workflows[i])
		if err != nil {
			return err
		}
	}

	p := models.Project{
		Key:         "TEST",
		Name:        "Test Project",
		CreatedDate: time.Now(),
		Lead:        "testadmin",
		TicketTypes: []string{
			"Epic",
			"Story",
			"Bug",
			"Feature Request",
		},

		Public:      true,
		FieldScheme: fs.ID,

		WorkflowScheme: []models.WorkflowMapping{
			{
				TicketType: "",
				Workflow:   workflows[0].ID,
			},
		},
	}

	p1 := p
	p1.Public = false
	p1.Key = "TEST2"
	p1.Name = "Another TEST Project"

	perms := map[models.Role][]permission.Permission{
		"Administrator": []permission.Permission{
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
		"Contributor": []permission.Permission{"VIEW_PROJECT",
			"CREATE_TICKET",
			"COMMENT_TICKET",
			"REMOVE_OWN_COMMENT",
			"EDIT_OWN_COMMENT",
			"TRANSITION_TICKET",
			"EDIT_TICKET",
		},
		"User": []permission.Permission{
			"VIEW_PROJECT",
			"CREATE_TICKET",
			"COMMENT_TICKET",
		},
	}

	for k, v := range perms {
		for _, prm := range v {
			roleMapping := models.RolePermission{
				Role:       k,
				Permission: prm,
			}

			p.Permissions = append(p.Permissions, roleMapping)
			p1.Permissions = append(p1.Permissions, roleMapping)
		}
	}

	p, err = r.Projects().Create(u1, p)
	if err != nil {
		return err
	}

	p1, err = r.Projects().Create(u1, p1)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		t := models.Ticket{
			Summary: "This is test ticket #" + strconv.Itoa(i),
			Description: `# Refugam in se fuit quae

## Pariter vel sine frustra

Lorem markdownum Diomede quid, ab oracula diligit; aut qui nam. Dum postquam tu
fecit *numerare dederat es* animae dederat, quem soror. Venae potentem minacia
summa precantem statque procubuisse et sui et deus sceleri?

1. Irascitur inter de cunctae arva tenet pectore
2. Tabo messibus
3. Duobus undae

## Truncis sulcat Stymphalide

Sollertius nomina plectrumque nec nec animos, Rhadamanthon figitur vulgata
hominum ad. Vulnere pendentemque soror incubuit lenta vertunt. Deae cepit
quotiensque toto Aenea curvamine cum non sua divus audet patriae si et fit
vineta. Aquas nimium: postquam hominum promissa!

    if (isdn >= personal_executable(cJquery)) {
        redundancy_firmware_guid = infringement;
        keystroke += pum_document(page_wins, icq_nanometer_malware +
                barInternal);
        mcaQueryMarketing(portLeak, guiPhreaking, thunderbolt(4, twainAtaLink));
    }
    addressTorrent = boot_character_website(linkedinVaporware, plugRightBoot);
    var megabit_standalone_of = nocSo + program_mouse + 26;

## Nostra est perdix annos et quas

Vellentem quaerit est umeros celsior navis intrat
[saepe](http://minosiuvenis.net/numen.html). Saxo vocet turris Athamanta
membris, semesaque: nate leto summos instabiles primosque avertite nostras tu
quies in [avidisque](http://www.templaaequora.net/). Summa se expulit perfide
mirum, suo brevi absentem umerus vultumque cognata. Nempe ipsi quod procul
verba, frusta, sed gemitu non huius odit; non aprica pedumque Hectoris, taxo.
Mentis vivit tori erubuit, qui flebile natura Echo percussis pallet?

- Ministros tumebat famuli
- Aristas per blandis
- Corpora qua Medea acu potentia inrita

Non Cipe reges, laetitiam filius sceleratum naidas, fortunaque occidit. Laeva et
ipsa divite, est ille ver verba vicisse, exsiliantque aprica illius, rapta?`,
			Reporter: users[rand.Intn(2)].Username,
			Assignee: users[rand.Intn(2)].Username,
			Type:     p.TicketTypes[rand.Intn(3)],
			Project:  p.Key,
		}

		t, err = r.Tickets().Create(u1, t)
		if err != nil {
			return err
		}

		for i := 0; i < rand.Intn(50); i++ {
			c := models.Comment{
				Author: users[rand.Intn(2)].Username,
				Body: `# Yo Dawg

I heard you like **markdown**.

So I put markdown in your comment.`,
			}

			_, err = r.Tickets().AddComment(u1, t.Key, c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}