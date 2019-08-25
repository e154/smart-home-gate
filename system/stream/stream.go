package stream

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"net/http"
	"strconv"
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
	adaptors *adaptors.Adaptors
	Hub      *Hub
}

func NewStreamService(hub *Hub,
	adaptors *adaptors.Adaptors) *StreamService {
	return &StreamService{
		Hub:      hub,
		adaptors: adaptors,
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

	accessToken := ctx.Request.Header.Get("X-API-Key")
	clientType := strings.ToLower(ctx.Request.Header.Get("X-Client-Type"))

	var token string

	var clientId string
	if accessToken != "" {
		data := strings.Split(accessToken, "-")
		if len(data) != 4 {
			ctx.AbortWithError(401, errors.New("unauthorized access"))
			return
		}

		requestRandomId := data[1]
		hash := data[3]

		timestamp, errw := strconv.Atoi(data[2])
		if errw != nil {
			ctx.AbortWithError(401, errors.New("unauthorized access"))
			return
		}

		if len(requestRandomId) < 8 {
			ctx.AbortWithError(401, errors.New("unauthorized access"))
			return
		}

		switch clientType {
		case ClientTypeServer:
			serverObj, err := w.adaptors.Server.GetById(data[0])
			if err != nil {
				ctx.AbortWithError(401, errors.New("unauthorized access"))
				return
			}

			clientId = serverObj.Id
			token = serverObj.Token

		case ClientTypeMobile:
			mobileObj, err := w.adaptors.Mobile.GetById(data[0])
			if err != nil {
				ctx.AbortWithError(401, errors.New("unauthorized access"))
				return
			}

			clientId = mobileObj.Id
			token = mobileObj.Token.String()
		default:
			ctx.AbortWithError(400, fmt.Errorf("unknown client type %v", clientType))
			return
		}

		if hash != common.Sha256(requestRandomId+token+fmt.Sprintf("%d", timestamp)) {
			ctx.AbortWithError(401, errors.New("unauthorized access"))
			return
		}

	} else {

		switch clientType {
		case ClientTypeServer:
		case ClientTypeMobile:
		default:
			ctx.AbortWithError(400, fmt.Errorf("unknown client type %v", clientType))
			return
		}
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
		Id:      clientId,
		Connect: conn,
		Ip:      ctx.ClientIP(),
		Send:    make(chan []byte),
		Token:   token,
		Type:    clientType,
	}

	go client.WritePump()
	w.Hub.AddClient(client)
}

func (s *StreamService) GetClientByIdAndType(clientId,
	clientType string) (client *Client, err error) {
	client, err = s.Hub.GetClientByIdAndType(clientId, clientType)
	return
}
