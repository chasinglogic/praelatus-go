// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
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

func handleHookEvent() {
	outEvents := make(chan hookEvent)

	for i := 0; i < webWorkers; i++ {
		go webWorker(outEvents)
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

func webWorker(inEvent chan hookEvent) {
	for {
		e := <-inEvent
		err := runHook(e.ticket, e.transition, e.hook)
		if err != nil {
			eventLog.Println("|Web Worker|", err)
		}
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

func runHook(ticket models.Ticket, transition models.Transition, hook models.Hook) error {
	body, err := renderBody(ticket, transition, hook)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(hook.Method, hook.Endpoint, body)
	if err != nil {
		return fmt.Errorf("error creating request %s: %s %s",
			ticket.Key, transition.Name, err.Error())
	}

	client := http.Client{}
	resp, err := client.Do(r)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if err != nil {
		return fmt.Errorf("request failed Ticket=%s: Status=%s Error=%s",
			ticket.Key, transition.Name, err.Error())
	}

	return nil
}
