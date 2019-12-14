package stream

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

const (
	ClientTypeServer = "server"
	ClientTypeMobile = "mobile"
)

type Client struct {
	Id      string
	Connect *websocket.Conn
	Ip      string
	Send    chan []byte
	Token   string
	Type    string
}

func (c *Client) Notify(t, b string) {

	msg, _ := json.Marshal(&map[string]interface{}{"type": "notify", "value": &map[string]interface{}{"type": t, "body": b}})

	c.Send <- msg

}

func (c *Client) Write(payload []byte) (err error) {
	c.Send <- payload
	return nil
}

func (c *Client) write(opCode int, payload []byte) error {
	c.Connect.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Connect.WriteMessage(opCode, payload)
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
