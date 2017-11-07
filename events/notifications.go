package events

import (
	"time"

	"github.com/praelatus/praelatus/events/event"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
)

var (
	notificationsEventChan = make(chan event.Event)

	// sendNotificationsChan = make(chan models.Notification)
	// notificationWorkers   = 10
)

func recordNofiticationEvent() {
	// for i := 0; i < notificationWorkers; i++ {
	// 	go sendNotificationWorker(result)
	// }

	for {
		e := <-notificationsEventChan

		n := models.Notification{
			Type:           string(e.Type()),
			ActionedTicket: e.Ticket().Key,
			ActioningUser:  e.ActioningUser().Username,
			Project:        e.Project().Key,
			CreatedDate:    time.Now(),
			Body:           e.String(),
			Read:           false,
		}

		eventLog.Println(n)

		for _, w := range e.Ticket().Watchers {
			if w == e.ActioningUser().Username {
				n.Watcher = w
			}

			_, err := repo.Notifications().Create(nil, n)
			if err != nil {
				eventLog.Println("|Notification Recorder|", err)
			}
		}
	}
}

// TODO: Support email notifications via this function
// func sendNotificationWorker(result chan Result) {

// }
