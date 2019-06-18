package stream

import (
	"encoding/json"
	"fmt"
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
	sessions    map[*Client]bool
	subscribers map[string]func(client *Client, value interface{})
	sync.Mutex
	broadcast chan []byte
	interrupt chan os.Signal
}

func NewHub() *Hub {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	hub := &Hub{
		sessions:    make(map[*Client]bool),
		broadcast:   make(chan []byte, maxMessageSize),
		subscribers: make(map[string]func(client *Client, value interface{})),
		interrupt:   interrupt,
	}
	go hub.Run()

	return hub
}

func (h *Hub) AddClient(client *Client) {

	defer func() {
		delete(h.sessions, client)
		log.Infof("websocket session from ip: %s closed", client.Ip)
	}()

	h.sessions[client] = true

	log.Infof("new websocket xsession established, from ip: %s", client.Ip)

	_ = client.Connect.SetReadDeadline(time.Now().Add(pongWait))
	client.Connect.SetPongHandler(func(string) error {
		_ = client.Connect.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		op, r, err := client.Connect.NextReader()
		if err != nil {
			log.Error(err.Error())
			break
		}
		switch op {
		case websocket.TextMessage:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				break
			}
			h.Recv(client, message)
		}
	}
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

func (h *Hub) Recv(client *Client, message []byte) {

	fmt.Printf("client(%v), message(%v)\n", client, string(message))

	re := map[string]interface{}{}
	if err := json.Unmarshal(message, &re); err != nil {
		log.Error(err.Error())
		return
	}

	for key, value := range re {

		switch key {
		//case "client_info":
			//client.UpdateInfo(value)

		default:
			for command, f := range h.subscribers {
				if key == command {
					f(client, value)
				}
			}
		}
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

func (h *Hub) Subscribe(command string, f func(client *Client, value interface{})) {
	log.Infof("subscribe %s", command)
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
}

func (h *Hub) UnSubscribe(command string) {
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
}
