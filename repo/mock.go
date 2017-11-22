// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// +build !release

package repo

import (
	"math/rand"
	"strconv"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/ql/ast"
	"gopkg.in/mgo.v2/bson"
)

var tickets []models.Ticket

func init() {
	for i := 0; i < 100; i++ {
		t := models.Ticket{
			Key:     "TEST-" + strconv.Itoa(i+1),
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

		for i := 0; i < rand.Intn(50); i++ {
			c := models.Comment{
				Author: users[rand.Intn(2)].Username,
				Body: `# Yo Dawg

I heard you like **markdown**.

So I put markdown in your comment.`,
			}

			t.Comments = append(t.Comments, c)
		}

		tickets = append(tickets, t)
	}

}

type mockRepo struct{}

func NewMockRepo() Repo {
	return mockRepo{}
}

type mockProjectRepo struct{}

func (pr mockProjectRepo) Get(u *models.User, uid string) (models.Project, error) {
	return p1, nil
}

func (pr mockProjectRepo) Search(u *models.User, query string) ([]models.Project, error) {
	return []models.Project{p, p1}, nil
}

func (pr mockProjectRepo) HasLead(u *models.User, lead models.User) ([]models.Project, error) {
	return []models.Project{p, p1}, nil
}

func (pr mockProjectRepo) Update(u *models.User, uid string, updated models.Project) error {
	return nil
}

func (pr mockProjectRepo) Create(u *models.User, project models.Project) (models.Project, error) {
	return project, nil
}

func (pr mockProjectRepo) Delete(u *models.User, uid string) error {
	return nil
}

type mockTicketRepo struct{}

func (t mockTicketRepo) Get(u *models.User, uid string) (models.Ticket, error) {
	return tickets[0], nil
}

func (t mockTicketRepo) Search(u *models.User, query ast.AST) ([]models.Ticket, error) {
	return tickets, nil
}

func (t mockTicketRepo) Update(u *models.User, uid string, updated models.Ticket) error {
	return nil
}

func (t mockTicketRepo) Create(u *models.User, ticket models.Ticket) (models.Ticket, error) {
	return ticket, nil
}

func (t mockTicketRepo) Delete(u *models.User, uid string) error {
	return nil
}

func (t mockTicketRepo) LabelSearch(u *models.User, query string) ([]string, error) {
	return []string{"label1", "label2"}, nil
}

func (t mockTicketRepo) AddComment(u *models.User, uid string, comment models.Comment) (models.Ticket, error) {
	tickets[0].Comments = append(tickets[0].Comments, comment)
	return tickets[0], nil
}

func (t mockTicketRepo) NextTicketKey(u *models.User, projectKey string) (string, error) {
	return projectKey + string(len(tickets)+1), nil
}

type mockUserRepo struct{}

func (ur mockUserRepo) Get(u *models.User, uid string) (models.User, error) {
	return *u1, nil
}

func (ur mockUserRepo) Search(u *models.User, query string) ([]models.User, error) {
	return users, nil
}

func (ur mockUserRepo) Update(u *models.User, uid string, updated models.User) error {
	return nil
}

func (ur mockUserRepo) Create(u *models.User, user models.User) (models.User, error) {
	return user, nil
}

func (ur mockUserRepo) Delete(u *models.User, uid string) error {
	return nil
}

type mockFieldRepo struct{}

func (fsr mockFieldRepo) Get(u *models.User, uid string) (models.FieldScheme, error) {
	// Hardcode to the ID expected in tests.
	fs.ID = "59e3f2026791c08e74da1bb2"
	return fs, nil
}

func (fsr mockFieldRepo) Search(u *models.User, query string) ([]models.FieldScheme, error) {
	return []models.FieldScheme{fs}, nil
}

func (fsr mockFieldRepo) Update(u *models.User, uid string, updated models.FieldScheme) error {
	return nil
}

func (fsr mockFieldRepo) Create(u *models.User, fieldScheme models.FieldScheme) (models.FieldScheme, error) {
	fieldScheme.ID = bson.NewObjectId()
	return fieldScheme, nil
}

func (fsr mockFieldRepo) Delete(u *models.User, uid string) error {
	return nil
}

type mockWorkflowRepo struct{}

func (wr mockWorkflowRepo) Get(u *models.User, uid string) (models.Workflow, error) {
	wrk := workflows[0]
	// Hardcode to the ID expected in tests.
	wrk.ID = "59e3f2026791c08e74da1bb2"
	return wrk, nil
}

func (wr mockWorkflowRepo) Search(u *models.User, query string) ([]models.Workflow, error) {
	return workflows, nil
}

func (wr mockWorkflowRepo) Update(u *models.User, uid string, updated models.Workflow) error {
	return nil
}

func (wr mockWorkflowRepo) Create(u *models.User, workflow models.Workflow) (models.Workflow, error) {
	workflow.ID = bson.NewObjectId()
	return workflow, nil
}

func (wr mockWorkflowRepo) Delete(u *models.User, uid string) error {
	return nil
}

type mockNotificationRepo struct{}

func (nr mockNotificationRepo) Create(u *models.User, notification models.Notification) (models.Notification, error) {
	notification.ID = bson.NewObjectId()
	return notification, nil
}

func (nr mockNotificationRepo) MarkRead(u *models.User, uid string) error {
	return nil
}

func (nr mockNotificationRepo) ForProject(u *models.User, project models.Project, onlyUnread bool, last int) ([]models.Notification, error) {
	return nil, nil
}

func (nr mockNotificationRepo) ForUser(u *models.User, user models.User, onlyUnread bool, last int) ([]models.Notification, error) {
	return nil, nil
}

func (nr mockNotificationRepo) ActivityForUser(u *models.User, user models.User, onlyUnread bool, last int) ([]models.Notification, error) {
	return nil, nil
}

func (m mockRepo) Projects() ProjectRepo {
	return mockProjectRepo{}
}

func (m mockRepo) Tickets() TicketRepo {
	return mockTicketRepo{}
}

func (m mockRepo) Users() UserRepo {
	return mockUserRepo{}
}

func (m mockRepo) Fields() FieldSchemeRepo {
	return mockFieldRepo{}
}

func (m mockRepo) Workflows() WorkflowRepo {
	return mockWorkflowRepo{}
}

func (m mockRepo) Notifications() NotificationRepo {
	return mockNotificationRepo{}
}

func (m mockRepo) Clean() error { return nil }
func (m mockRepo) Test() error  { return nil }
func (m mockRepo) Init() error  { return nil }
