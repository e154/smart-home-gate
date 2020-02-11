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
	"github.com/e154/smart-home-gate/adaptors"
	m "github.com/e154/smart-home-gate/models"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

type Hub struct {
	adaptors        *adaptors.Adaptors
	broadcast       chan []byte
	interrupt       chan os.Signal
	sessionsLock    sync.Mutex
	sessions        map[*Client]bool
	subscribersLock sync.Mutex
	subscribers     map[string]func(client *Client, msg Message)
}

func NewHub(adaptors *adaptors.Adaptors) *Hub {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	hub := &Hub{
		adaptors:    adaptors,
		sessions:    make(map[*Client]bool),
		broadcast:   make(chan []byte, maxMessageSize),
		subscribers: make(map[string]func(client *Client, msg Message)),
		interrupt:   interrupt,
	}
	go hub.Run()

	return hub
}

func (h *Hub) AddClient(client *Client) {

	clientId := client.Id
	if clientId == "" {
		clientId = "Empty"
	}

	defer func() {
		h.sessionsLock.Lock()
		if ok := h.sessions[client]; ok {
			delete(h.sessions, client)
		}
		h.sessionsLock.Unlock()
		log.Infof("websocket session from ip(%s) closed, id(%s)", client.Ip, clientId)
	}()

	h.sessionsLock.Lock()
	h.sessions[client] = true
	h.sessionsLock.Unlock()

	log.Infof("new websocket session established, from ip(%s), type(%v), id(%s)", client.Ip, client.Type, clientId)

	_ = client.Connect.SetReadDeadline(time.Now().Add(pongWait))
	client.Connect.SetPongHandler(func(string) error {
		_ = client.Connect.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		op, r, err := client.Connect.NextReader()
		if err != nil {
			log.Debug(err.Error())
			break
		}
		switch op {
		case websocket.TextMessage:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				break
			}
			h.Recv(client, message)
		default:
			log.Warningf("unknown message type(%v)", op)
		}
	}

	log.Infof("websocket session closed id(%s)", client.Id)
}

func (h *Hub) GetClientByIdAndType(clientId, clientType string) (client *Client, err error) {

	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()
	
	for cli, _ := range h.sessions {
		if cli.Id == clientId && cli.Type == clientType {
			client = cli
			return
		}
	}

	return
}

func (h *Hub) Run() {

	for {
		select {
		case m := <-h.broadcast:
			h.sessionsLock.Lock()
			for client := range h.sessions {
				client.Send <- m
			}
			h.sessionsLock.Unlock()
		case <-h.interrupt:
			//fmt.Println("Close websocket client session")
			h.sessionsLock.Lock()
			for client := range h.sessions {
				client.Close()
				delete(h.sessions, client)
			}
			h.sessionsLock.Unlock()
		}

	}
}

//
// client --> server
//
// server --> [client1, client2]
//
func (h *Hub) Recv(client *Client, b []byte) {

	clientId := client.Id
	if clientId == "" {
		clientId = "Empty"
	}
	log.Debugf("Receive message from client type(%v), clientId(%v)", client.Type, clientId)

	//fmt.Printf("client(%v), message(%v)\n", client, string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch client.Type {
	case ClientTypeServer:
		switch msg.Command {
		default:
			if f := h.GetCommandFromSubscribers(msg.Command); f != nil {
				f(client, msg)
				return
			}
		}
		var server *m.Server
		if server, err = h.adaptors.Server.GetById(client.Id); err != nil {
			log.Error(err.Error())
			return
		}
		if server.Mobiles == nil && len(server.Mobiles) == 0 {
			log.Info("clients for message not found")
			return
		}
		h.sessionsLock.Lock()
		for client, _ := range h.sessions {
			if client.Type == ClientTypeServer {
				continue
			}
			for _, mobile := range server.Mobiles {
				if mobile.Id == client.Id {
					log.Debugf("Resend to mobile client: %v", client.Id)
					client.Send <- b
				}
			}
		}
		h.sessionsLock.Unlock()

	case ClientTypeMobile:
		var mobile *m.Mobile
		if mobile, err = h.adaptors.Mobile.GetById(client.Id); err != nil {
			log.Error(err.Error())
			return
		}
		h.sessionsLock.Lock()
		for client, _ := range h.sessions {
			if client.Type == ClientTypeMobile {
				continue
			}
			if mobile.ServerId == client.Id {
				log.Debugf("Resend to server: %v", client.Id)
				client.Send <- b
				h.sessionsLock.Unlock()
				return
			}
			log.Warningf("server %s not found", mobile.ServerId)
		}
		h.sessionsLock.Unlock()
	default:
		log.Errorf("unknown client type: %v", client.Type)
	}
}

func (h *Hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Clients() (clients []*Client) {
	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()
	clients = []*Client{}
	for c := range h.sessions {
		clients = append(clients, c)
	}

	return
}

func (h *Hub) Subscribe(command string, f func(client *Client, msg Message)) {
	h.subscribersLock.Lock()
	defer h.subscribersLock.Unlock()
	log.Infof("subscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
}

func (h *Hub) UnSubscribe(command string) {
	h.subscribersLock.Lock()
	defer h.subscribersLock.Unlock()
	log.Infof("unsubscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
}

func (h *Hub) GetCommandFromSubscribers(cmd string) func(client *Client, msg Message) {
	h.subscribersLock.Lock()
	defer h.subscribersLock.Unlock()
	for command, f := range h.subscribers {
		if cmd == command {
			return f
		}
	}
	return nil
}
