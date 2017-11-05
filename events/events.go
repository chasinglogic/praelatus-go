// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
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

// evm is the global event manager which is interfaced with using the functions
// whose names match the methods of an EventManager
var evm = &EventManager{
	Result: make(chan Result, 10),
	Listeners: []chan event.Event{
		hookEventChan,
	},
}

var eventLog = log.New(config.LogWriter(), "EVENT", log.LstdFlags)

// Stop is a global channel used to stop all running event managers
var Stop = make(chan int)

// ResultChan returns the result channel used on the global EventManager
func ResultChan() chan Result {
	return evm.Result
}

// Run starts the global event manager, should be called in a go routine normally
func Run() {
	go handleHookEvent(evm.Result)

	for {
		select {
		case res := <-evm.Result:
			if res.Error != nil {
				eventLog.Printf("handler %s failed with error %s\n",
					res.Reporter, res.Error.Error())
			}
		case <-Stop:
			return
		}
	}
}

// FireEvent calls the method of the same name on the global EventManager
func FireEvent(e event.Event) {
	evm.FireEvent(e)
}

// Result contains metadata sent back by an event handler
type Result struct {
	Reporter string
	Error    error
}

// EventManager is the fan-in point for all available event listeners
type EventManager struct {
	Result    chan Result
	Listeners []chan event.Event
}

// FireEvent sends the given event to all registered listeners of the given
// event manager
func (e *EventManager) FireEvent(ev event.Event) {
	for _, listener := range e.Listeners {
		listener <- ev
	}
	// Return to guarantee that this is killed if called in a goroutine
	return
}

// RegisterListener adds a listener to the EventManager
func (e *EventManager) RegisterListener(ev chan event.Event) {
	e.Listeners = append(e.Listeners, ev)
}

// New will return a blank event manager handling allocations for you as
// appropriate
func New() *EventManager {
	return &EventManager{
		Result:    make(chan Result),
		Listeners: []chan event.Event{},
	}
}
