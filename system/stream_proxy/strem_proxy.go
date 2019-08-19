package stream_proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/common"
	"github.com/e154/smart-home-gate/common/debug"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/e154/smart-home-gate/system/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

var (
	log = logging.MustGetLogger("stream_proxy")
)

type StreamProxy struct {
	engine        *gin.Engine
	streamService *stream.StreamService
	adaptors      *adaptors.Adaptors
}

func NewStreamProxy(httpsServer *server.Server,
	streamService *stream.StreamService,
	graceful *graceful_service.GracefulService,
	adaptors *adaptors.Adaptors) (proxy *StreamProxy) {
	proxy = &StreamProxy{
		engine:        httpsServer.GetEngine(),
		streamService: streamService,
		adaptors:      adaptors,
	}

	graceful.Subscribe(proxy)

	return
}

func (s *StreamProxy) Start() {
	log.Info("start stream proxy")
	s.streamService.Subscribe("chanel.server", s.DoAction)
	s.engine.Any("/server/*any", s.auth, s.controller)
}

func (s *StreamProxy) Shutdown() {
	s.streamService.UnSubscribe("chanel.server")
	return
}

func (s *StreamProxy) DoAction(client *stream.Client, message stream.Message) {

	fmt.Println("------")
	debug.Println(client)
	fmt.Println("------")
	debug.Println(message)

	return
}

func (s *StreamProxy) execRequest() {

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/", nil)
	request.SetBasicAuth("admin", "admin")

	s.engine.ServeHTTP(recorder, request)
	fmt.Println(recorder.Code)
	fmt.Println(recorder.Body)
}

// access_token
func (s *StreamProxy) getToken(ctx *gin.Context) (accessToken string, err error) {

	if accessToken = ctx.Request.Header.Get("server_access_token"); accessToken != "" {
		return
	}

	if accessToken = ctx.Request.Header.Get("ServerAuthorization"); accessToken != "" {
		return
	}

	if accessToken = ctx.Request.URL.Query().Get("server_access_token"); accessToken != "" {
		return
	}

	return
}

func (s *StreamProxy) auth(ctx *gin.Context) {

	var err error

	// get access_token
	var accessToken string
	if accessToken, err = s.getToken(ctx); err != nil || accessToken == "" {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	data := strings.Split(accessToken, "-")
	if len(data) != 4 {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	mobileClientId := data[0]
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

	var mobileObj *m.Mobile
	if mobileObj, err = s.adaptors.Mobile.GetById(mobileClientId); err != nil {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	var serverObj *m.Server
	if serverObj, err = s.adaptors.Server.GetById(mobileObj.ServerId); err != nil {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	if hash != common.Sha256(requestRandomId+mobileObj.Token.String()+fmt.Sprintf("%d", timestamp)) {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	if serverObj != nil {
		if ctx.Keys == nil {
			ctx.Keys = make(map[string]interface{})
		}
		ctx.Keys["server"] = serverObj
		return
	}

	log.Warningf(fmt.Sprintf("access denied token: %s", accessToken))

	ctx.AbortWithError(403, errors.New("unauthorized access"))
}

func (s *StreamProxy) controller(ctx *gin.Context) {

	var serverObj *m.Server
	if _, ok := ctx.Keys["server"]; ok {
		serverObj = ctx.Keys["server"].(*m.Server)
	} else {
		ctx.String(http.StatusNotFound, "server not found")
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	streamRequestModel := &StreamRequestModel{
		URI:    ctx.Request.RequestURI,
		Method: strings.ToUpper(ctx.Request.Method),
		Body:   body,
		Header: ctx.Request.Header,
	}

	fmt.Printf("serverId: %v\n", serverObj.Id)
	fmt.Printf("streamRequestModel: %v\n", streamRequestModel)

	var client *stream.Client
	if client, err = s.streamService.GetClientByIdAndType(serverObj.Id, stream.ClientTypeServer); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if client == nil {
		ctx.String(http.StatusNotFound, "server not found")
		return
	}

	debug.Println(client)

	payload := map[string]interface{}{
		"request": streamRequestModel,
	}

	s.Send("mobile_gate_proxy", payload, client, ctx, func(msg stream.Message) {
		fmt.Println("ok")
		debug.Println(msg)
	})

	return
}

func (g *StreamProxy) Send(command string, payload map[string]interface{}, client *stream.Client, ctx *gin.Context, f func(msg stream.Message)) (err error) {

	done := make(chan struct{})

	message := stream.Message{
		Id:      uuid.NewV4(),
		Command: command,
		Payload: payload,
	}

	g.streamService.Subscribe(message.Id.String(), func(client *stream.Client, msg stream.Message) {
		g.streamService.UnSubscribe(message.Id.String())
		f(msg)
		done <- struct{}{}
	})

	msg, _ := json.Marshal(message)
	if err := client.Write(websocket.TextMessage, msg); err != nil {
		log.Error(err.Error())
	}

	select {
	case <-time.After(2 * time.Second):
		ctx.AbortWithStatus(http.StatusRequestTimeout)
	case <-done:
	case <-ctx.Done():
		ctx.AbortWithStatus(http.StatusRequestTimeout)
	}

	return
}
