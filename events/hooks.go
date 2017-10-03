// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package events

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
)

var (
	hookEventChan = make(chan models.Event)
	webWorkers    = config.WebWorkers()
)

type hookEvent struct {
	ticket     models.Ticket
	transition models.Transition
	hook       models.Hook
}

func handleHookEvent(result chan Result) {
	outEvents := make(chan hookEvent)

	for i := 0; i < webWorkers; i++ {
		go sendRequest(outEvents, result)
	}

	for {
		event := <-hookEventChan

		transition, ok := event.Data.(models.Transition)
		if !ok {
			continue
		}

		for _, hook := range transition.Hooks {
			outEvents <- hookEvent{
				ticket:     event.Ticket,
				transition: transition,
				hook:       hook,
			}
		}
	}
}

func sendRequest(events chan hookEvent, result chan Result) {
	for {
		event := <-events

		res := Result{Reporter: "Hook Handler", Success: true}

		tmpl, err := template.New("hook-body").Parse(event.hook.Body)
		if err != nil {
			e := fmt.Sprintf("Error parsing body %s: %s %s\n",
				event.ticket.Key, event.transition.Name, err.Error())
			res.Success = false
			res.Error = errors.New(e)
			result <- res
			continue
		}

		body := bytes.NewBuffer([]byte{})

		err = tmpl.Execute(body, event.ticket)
		if err != nil {
			e := fmt.Sprintf("Error rendering body %s: %s %s\n",
				event.ticket.Key, event.transition.Name, err.Error())
			res.Success = false
			res.Error = errors.New(e)
			result <- res
			continue
		}

		r, err := http.NewRequest(event.hook.Method, event.hook.Endpoint, body)
		if err != nil {
			e := fmt.Sprintf("Error creating request %s: %s %s\n",
				event.ticket.Key, event.transition.Name, err.Error())
			res.Success = false
			res.Error = errors.New(e)
			result <- res
			continue
		}

		client := http.Client{}

		resp, err := client.Do(r)

		if resp != nil {
			_ = resp.Body.Close()
		}

		if err != nil {
			e := fmt.Sprintf("REQUEST FAILED Ticket=%s: Status=%s Error=%s\n",
				event.ticket.Key, event.transition.Name, err.Error())
			res.Success = false
			res.Error = errors.New(e)
			result <- res
			continue
		}

		result <- res
	}
}
