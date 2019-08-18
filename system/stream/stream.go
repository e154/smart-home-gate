package stream

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"net/http"
	"strings"
)

var (
	log        = logging.MustGetLogger("stream")
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type StreamService struct {
	Hub *Hub
}

func NewStreamService(hub *Hub) *StreamService {
	return &StreamService{
		Hub: hub,
	}
}

func (s *StreamService) Broadcast(message []byte) {
	s.Hub.Broadcast(message)
}

func (s *StreamService) Subscribe(command string, f func(client *Client, msg Message)) {
	s.Hub.Subscribe(command, f)
}

func (s *StreamService) UnSubscribe(command string) {
	s.Hub.UnSubscribe(command)
}

func (w *StreamService) Ws(ctx *gin.Context) {

	clientType := strings.ToLower(ctx.Request.Header.Get("X-Client-Type"))
	switch clientType {
	case ClientTypeServer, ClientTypeMobile:
	default:
		ctx.AbortWithError(400, fmt.Errorf("unknown client type %v", clientType))
		return
	}

	// CORS
	ctx.Writer.Header().Del("Access-Control-Allow-Credentials")

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("Failed to set websocket upgrade: %v", err)
		return
	}
	if _, ok := err.(websocket.HandshakeError); ok {
		ctx.AbortWithError(400, errors.New("not a websocket handshake"))
		return
	}

	client := &Client{
		Connect: conn,
		Ip:      ctx.ClientIP(),
		Send:    make(chan []byte),
		Token:   ctx.Request.Header.Get("X-API-Key"),
		Type:    clientType,
	}

	go client.WritePump()
	w.Hub.AddClient(client)
}

func (s *StreamService) GetClientByToken(token string) (client *Client, err error) {
	client, err = s.Hub.GetClientByToken(token)
	return
}
