package events

import (
	"testing"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

var r = repo.NewMockRepo()

func TestHookListener(t *testing.T) {
	evm := New()
	evm.RegisterListener(hookEventChan)

	go handleHookEvent(evm.Result)

	tk, _ := r.Tickets().Get(nil, "")

	e := models.Event{
		Ticket: tk,
		Data: models.Transition{
			Name:     "In Progress",
			ToStatus: "Done",
			Hooks: []models.Hook{
				{
					Method:   "GET",
					Endpoint: "https://google.com",
					Body:     "",
				},
			},
		},
	}

	hookEventChan <- e

	res := <-evm.Result
	if !res.Success {
		t.Errorf("%s: %s\n", res.Reporter, res.Error.Error())
	}
}
