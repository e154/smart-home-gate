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
	adaptors    *adaptors.Adaptors
	sessions    map[*Client]bool
	subscribers map[string]func(client *Client, msg Message)
	sync.Mutex
	broadcast chan []byte
	interrupt chan os.Signal
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

	defer func() {
		if ok := h.sessions[client]; ok {
			delete(h.sessions, client)
		}
		log.Infof("websocket session from ip(%s) closed, id(%s)", client.Ip, client.Id)
	}()

	h.sessions[client] = true

	log.Infof("new websocket session established, from ip(%s), type(%v), id(%s)", client.Ip, client.Type, client.Id)

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
			for client := range h.sessions {
				client.Send <- m
			}
		case <-h.interrupt:
			//fmt.Println("Close websocket client session")
			for client := range h.sessions {
				client.Close()
				delete(h.sessions, client)
			}
		}

	}
}

//
// client --> server
//
// server --> [client1, client2]
//
func (h *Hub) Recv(client *Client, b []byte) {

	log.Debugf("Receive message from client type(%v): %v", client.Type, client.Id)

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
			for command, f := range h.subscribers {

				if msg.Command == command {
					f(client, msg)
					return
				}
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

	case ClientTypeMobile:
		var mobile *m.Mobile
		if mobile, err = h.adaptors.Mobile.GetById(client.Id); err != nil {
			log.Error(err.Error())
			return
		}
		for client, _ := range h.sessions {
			if client.Type == ClientTypeMobile {
				continue
			}
			if mobile.ServerId == client.Id {
				log.Debugf("Resend to server: %v", client.Id)
				client.Send <- b
				return
			}
			log.Warningf("server %s not found", mobile.ServerId	)
		}
	default:
		log.Errorf("unknown client type: %v", client.Type)
	}
}

func (h *Hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *Hub) Broadcast(message []byte) {
	h.Lock()
	h.broadcast <- message
	h.Unlock()
}

func (h *Hub) Clients() (clients []*Client) {

	clients = []*Client{}
	for c := range h.sessions {
		clients = append(clients, c)
	}

	return
}

func (h *Hub) Subscribe(command string, f func(client *Client, msg Message)) {
	log.Infof("subscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
}

func (h *Hub) UnSubscribe(command string) {
	log.Infof("unsubscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
}
