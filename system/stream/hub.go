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
	"errors"
	"github.com/e154/smart-home-gate/adaptors"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/gorilla/websocket"
	"io/ioutil"
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
	interrupt       chan struct{}
	sessionsLock    sync.Mutex
	servers         map[string]*Client
	mobiles         map[string]*Client
	subscribersLock sync.Mutex
	subscribers     map[string]func(client *Client, msg Message)
}

func NewHub(adaptors *adaptors.Adaptors,
	graceful_service *graceful_service.GracefulService) *Hub {

	hub := &Hub{
		adaptors:    adaptors,
		servers:     make(map[string]*Client),
		mobiles:     make(map[string]*Client),
		subscribers: make(map[string]func(client *Client, msg Message)),
		interrupt:   make(chan struct{}, 1),
	}

	graceful_service.Subscribe(hub)

	go hub.Run()

	return hub
}

func (h *Hub) Shutdown() {
	h.interrupt <- struct{}{}
}

func (h *Hub) AddClient(client *Client) {

	clientId := client.Id
	if clientId == "" {
		clientId = "Empty"
	}

	defer func() {
		h.sessionsLock.Lock()
		if client.Type == ClientTypeServer {
			delete(h.servers, client.Id)
		} else {
			delete(h.mobiles, client.Id)
		}
		h.sessionsLock.Unlock()
		log.Infof("websocket session from ip(%s) closed, id(%s)", client.Ip, clientId)
	}()

	h.sessionsLock.Lock()
	if client.Type == ClientTypeServer {
		h.servers[client.Id] = client
	} else {
		h.mobiles[client.Id] = client
	}
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
		client.receivedInc()
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

func (h *Hub) GetClientServer(clientId string) (client *Client, err error) {

	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()

	var ok bool
	if client, ok = h.servers[clientId]; !ok {
		err = errors.New("not found")
	}

	return
}

func (h *Hub) Run() {

	for {
		select {
		case <-h.interrupt:
			h.sessionsLock.Lock()
			for _, client := range h.servers {
				client.Close()
				delete(h.servers, client.Id)
			}
			for _, mobile := range h.mobiles {
				mobile.Close()
				delete(h.servers, mobile.Id)
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

	lastMsgTime := client.getLastMsgTime()
	client.updateLastMsgTime()

	if client.Type == ClientTypeMobile && lastMsgTime < 0.1 {
		log.Warningf("Rejected message from client type(%v), clientId(%v), lastMsgTime(%v)", client.Type, clientId, lastMsgTime)
		return
	}

	//log.Debugf("Receive message from client type(%v), clientId(%v), lastMsgTime(%v)", client.Type, clientId, lastMsgTime)

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

		clients, err := h.GetServerClients(client)
		if err != nil {
			log.Error(err.Error())
		}
		for _, client := range clients {
			client.Send <- b
		}

	case ClientTypeMobile:
		if server, err := h.GetServer(client); err == nil {
			server.Send <- b
			return
		}
	default:
		log.Errorf("unknown client type: %v", client.Type)
	}
}

func (h *Hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *Hub) Subscribe(command string, f func(client *Client, msg Message)) {
	h.subscribersLock.Lock()
	defer h.subscribersLock.Unlock()
	//log.Infof("subscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
}

func (h *Hub) UnSubscribe(command string) {
	h.subscribersLock.Lock()
	defer h.subscribersLock.Unlock()
	//log.Infof("unsubscribe %s", command)
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

func (h *Hub) GetServer(cli *Client) (server *Client, err error) {

	var mobile *m.Mobile
	if mobile, err = h.adaptors.Mobile.GetById(cli.Id); err != nil {
		log.Error(err.Error())
		return
	}
	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()

	var ok bool
	if server, ok = h.servers[mobile.ServerId]; ok {
		log.Debugf("mobile --> server(%v)", mobile.ServerId)
	} else {
		log.Warningf("server %s not found", mobile.ServerId)
		err = errors.New("not found")
	}

	return
}

func (h *Hub) GetServerClients(cli *Client) (clients []*Client, err error) {

	var server *m.Server
	if server, err = h.adaptors.Server.GetById(cli.Id); err != nil {
		log.Error(err.Error())
		return
	}
	if server.Mobiles == nil && len(server.Mobiles) == 0 {
		log.Info("clients for message not found")
		return
	}
	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()

	for _, mobile := range server.Mobiles {
		if client, ok := h.mobiles[mobile.Id]; ok {
			log.Debugf("server(%v) --> mobile(%v)", server.Id, mobile.Id)
			clients = append(clients, client)
		}
	}

	return
}
