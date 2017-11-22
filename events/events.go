// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package events handles sending notifications to users and websockets
// when system events happen. Such as firing webhooks on workflow
// transitions.
package events

import (
	"log"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/events/event"
)

// TODO it's possible that one go routine taking a long time (such as the hook
// event manager thread) can lock up the whole event system, need to investigate
// the cost of spinning up new threads and channels to determine what the best
// way to go is.

var eventLog = log.New(config.LogWriter(), "[EVENT] ", log.LstdFlags)

// Stop is a global channel used to stop the running event managers
var Stop = make(chan int)

// EventManager manages registered event listeners and dispatching those events
type EventManager struct {
	Listeners []chan event.Event
}

// evm is the global event manager which is interfaced with using the functions
// whose names match the methods of an EventManager
var evm = EventManager{
	Listeners: []chan event.Event{
		hookEventChan,
		notificationsEventChan,
	},
}

// Run starts the global event manager, should be called in a go routine normally
func Run() {
	go handleHookEvent()
	go recordNofiticationEvent()

	for {
		_ = <-Stop
		eventLog.Println("Event Manager Shutting Down")
		return
	}
}

// FireEvent calls the method of the same name on the global EventManager
func (em EventManager) FireEvent(e event.Event) {
	eventLog.Println("fired", e.Type(), "for", e.Ticket().Key)
	// FireEvent sends the given event to all registered listeners
	for _, listener := range em.Listeners {
		listener <- e
		eventLog.Println("sent")
	}

	// Return to guarantee that this is killed if called in a goroutine
	return
}

// RegisterListener adds a listener to the EventManager
func (em EventManager) RegisterListener(ev chan event.Event) {
	em.Listeners = append(em.Listeners, ev)
}

// RegisterListener calls the method of the same name on the global EventManager
func RegisterListener(ev chan event.Event) {
	evm.RegisterListener(ev)
}

// FireEvent calls the method of the same name on the global EventManager
func FireEvent(ev event.Event) {
	evm.FireEvent(ev)
}
