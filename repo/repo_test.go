package repo

import (
	"math/rand"
	"strconv"

	"github.com/praelatus/praelatus/models"
)

var tickets []models.Ticket

func init() {
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

			t.Comments, err = append(t.Comments, c)
			if err != nil {
				return err
			}
		}
	}

}

type mockRepo struct{}

func NewMockRepo() Repo {
	return mockRepo{}
}

// TODO: Implement
type mockProjectRepo struct{}
type mockTicketRepo struct{}
type mockUserRepo struct{}
type mockFieldRepo struct{}
type mockWorkflowRepo struct{}

func (m mockRepo) Projects() ProjectRepo {
	return mockProjectRepo{}
}

func (m mockRepo) Tickets() TicketRepo {
	return mockTicketRepo{}
}

func (m mockRepo) Projects() ProjectRepo {
	return mockProjectRepo{}
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

func (m mockRepo) Clean() error { return nil }
func (m mockRepo) Test() error  { return nil }
