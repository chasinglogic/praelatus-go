package events

import (
	"github.com/praelatus/praelatus/events/event"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

var recordNotificationsEventChan = make(chan event.Event)
var sendNotificationsEventChan = make(chan event.Event)

func recordNofiticationEvent(result chan Result) {
	for {
		e := <-recordNotificationsEventChan

		n := models.Notification{
			Type:           string(e.Type()),
			ActionedTicket: e.Ticket().Key,
			ActioningUser:  e.ActioningUser().Username,
			Project:        e.Project().Key,
			Body:           e.String(),
			Read:           false,
		}

		for _, w := range e.Ticket().Watchers {
			n.Watcher = w.Username
			_, err := repo.Notifications().Create(nil, n)
			if err != nil {
				result <- Result{
					Name:  "Notification Recorder",
					Error: err,
				}
			}
		}
	}
}

func sendNotificationsEvent(result chan Result) {

}
