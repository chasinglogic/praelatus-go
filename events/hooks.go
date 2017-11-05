// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package events

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/events/event"
	"github.com/praelatus/praelatus/models"
)

var (
	hookEventChan = make(chan event.Event)
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
		go webWorker(outEvents, result)
	}

	for {
		e := <-hookEventChan

		if e.Type() != event.TransitionEvent {
			continue
		}

		transition, ok := e.Data().(models.Transition)
		if !ok {
			continue
		}

		for _, hook := range transition.Hooks {
			outEvents <- hookEvent{
				ticket:     e.Ticket(),
				hook:       hook,
				transition: transition,
			}
		}
	}
}

func webWorker(inEvent chan hookEvent, resultChan chan Result) {
	for {
		e := <-inEvent
		resultChan <- runHook(e.ticket, e.transition, e.hook)
	}
}

func renderBody(ticket models.Ticket, transition models.Transition, hook models.Hook) (io.Reader, error) {
	tmpl, err := template.New("hook-body").Parse(hook.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing body %s: %s %s",
			ticket.Key, transition.Name, err.Error())
	}

	body := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(body, ticket)
	if err != nil {
		return nil, fmt.Errorf("error rendering body %s: %s %s",
			ticket.Key, transition.Name, err.Error())
	}

	return body, nil
}

func runHook(ticket models.Ticket, transition models.Transition, hook models.Hook) Result {
	res := Result{Reporter: "Hook Handler: " + transition.Name}

	body, err := renderBody(ticket, transition, hook)
	if err != nil {
		res.Error = err
		return res
	}

	r, err := http.NewRequest(hook.Method, hook.Endpoint, body)
	if err != nil {
		res.Error = fmt.Errorf("error creating request %s: %s %s",
			ticket.Key, transition.Name, err.Error())
		return res
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if err != nil {
		res.Error = fmt.Errorf("request failed Ticket=%s: Status=%s Error=%s",
			ticket.Key, transition.Name, err.Error())
		return res
	}

	return res
}
