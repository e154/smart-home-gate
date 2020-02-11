// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package stream

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Client struct {
	Id          string
	Connect     *websocket.Conn
	Ip          string
	Send        chan []byte
	Token       string
	Type        ClientType
	writeLock   sync.Mutex
	lastMsgTime time.Time
	connected   time.Time
	server      *Client
}

func (c *Client) Notify(t, b string) {

	msg, _ := json.Marshal(&map[string]interface{}{"type": "notify", "value": &map[string]interface{}{"type": t, "body": b}})

	c.Send <- msg

}

func (c *Client) Write(payload []byte) (err error) {
	c.Send <- payload
	return nil
}

func (c *Client) write(opCode int, payload []byte) (err error) {
	c.writeLock.Lock()
	c.Connect.SetWriteDeadline(time.Now().Add(writeWait))
	err = c.Connect.WriteMessage(opCode, payload)
	c.writeLock.Unlock()
	return
}

// send message to client
//
func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		if c.Connect != nil {
			_ = c.Connect.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:

			if !ok {
				_ = c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) Close() {
	_ = c.write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *Client) getLastMsgTime() float64 {
	return time.Since(c.lastMsgTime).Seconds()
}

func (c *Client) updateLastMsgTime() {
	c.lastMsgTime = time.Now()
}

func (c *Client) UpdateServer(server *Client) {
	c.server = server
}
