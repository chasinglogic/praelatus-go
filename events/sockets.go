// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package events

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/praelatus/praelatus/events/event"
)

// TODO: Add this to event system need to have a better design first

var wsEventChan = make(chan event.Event, 10)

// ManagedWebsocket wraps a websocket connection and should be managed inside of a
// single go routine. Websockets are not safe for Concurrent reads and writes
// and so the ManagedWebsocket is used to guarantee that only one go routine is
// reading or writing at any given time.
type ManagedWebsocket struct {
	Socket *websocket.Conn
}

// WSManager is used to send events to listening websockets
type WSManager struct {
	activeWS []ManagedWebsocket
}

var wsm = WSManager{
	activeWS: make([]ManagedWebsocket, 0),
}

// AddWs will upgrade the given http connection to a websocket and register it
// with the Websocket Manager
func (e *WSManager) AddWs(w http.ResponseWriter, r *http.Request, h http.Header) error {
	conn, err := websocket.Upgrade(w, r, h, 1024, 1024)
	if err != nil {
		return err
	}

	ws := ManagedWebsocket{
		Socket: conn,
	}

	e.activeWS = append(e.activeWS, ws)
	return nil
}

// AddWs calls the method of the same name on the global WSManager
func AddWs(w http.ResponseWriter, r *http.Request, h http.Header) {
	wsm.AddWs(w, r, h)
}
