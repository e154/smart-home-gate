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

package stream_proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/common"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/e154/smart-home-gate/system/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
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

	proxy.engine.Any("/server/*any", proxy.auth, proxy.controller)

	graceful.Subscribe(proxy)

	return
}

func (s *StreamProxy) Start() {
	log.Info("start stream proxy")
	s.streamService.Subscribe("chanel.server", s.DoAction)
}

func (s *StreamProxy) Shutdown() {
	s.streamService.UnSubscribe("chanel.server")
	return
}

func (s *StreamProxy) DoAction(client *stream.Client, message stream.Message) {

	//fmt.Println("------")
	//debug.Println(client)
	//fmt.Println("------")
	//debug.Println(message)

	return
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

	accessToken = strings.TrimSpace(accessToken)

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

	defer func() {
		var err error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
		if err != nil {
			log.Error(err.Error())
		}
	}()

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

	url := strings.Replace(ctx.Request.RequestURI, "/server", "", -1)
	streamRequestModel := &StreamRequestModel{
		URI:    url,
		Method: strings.ToUpper(ctx.Request.Method),
		Body:   body,
		Header: ctx.Request.Header,
	}

	//fmt.Println("-------")
	//fmt.Println("request")
	//fmt.Println("-------")
	//fmt.Println(streamRequestModel.URI)
	//fmt.Println(streamRequestModel.Method)
	//fmt.Println(streamRequestModel.Header)
	//fmt.Println(string(streamRequestModel.Body))

	//fmt.Printf("serverId: %v\n", serverObj.Id)
	//fmt.Printf("streamRequestModel: %v\n", streamRequestModel)

	var client *stream.Client
	if client, err = s.streamService.GetClientByIdAndType(serverObj.Id, stream.ClientTypeServer); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if client == nil {
		ctx.String(http.StatusNotFound, "server not found")
		return
	}

	//debug.Println(client)

	payload := map[string]interface{}{
		"request": streamRequestModel,
	}

	err = s.Send("mobile_gate_proxy", payload, client, ctx, func(msg stream.Message) {

		//debug.Println(msg.Payload)

		if _, ok := msg.Payload["response"]; !ok {
			log.Error("no response field from payload")
			return
		}

		r := &StreamResponseModel{}
		if err = common.Copy(&r, msg.Payload["response"], common.JsonEngine); err != nil {
			log.Error(err.Error())
		}

		//fmt.Println("----------")
		//fmt.Println("response")
		//fmt.Println("----------")
		//fmt.Println(r.Code)
		//fmt.Println(r.Header)
		//fmt.Println(string(r.Body))

		for k, _ := range r.Header {
			ctx.Header(k, r.Header.Get(k))
		}

		if ctx.IsAborted() {
			log.Warning("context is aborted")
			return
		}

		ctx.Render(r.Code, render.Data{Data: r.Body})
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

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
	if err := client.Write(msg); err != nil {
		log.Error(err.Error())
	}

	select {
	case <-time.After(30 * time.Second):
		if !ctx.IsAborted() {
			ctx.AbortWithStatus(http.StatusRequestTimeout)
		}
	case <-done:
	case <-ctx.Done():
		if !ctx.IsAborted() {
			ctx.AbortWithStatus(http.StatusRequestTimeout)
		}
	}

	return
}
