package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home-gate/system/stream"
	"reflect"
)

type ControllerServer struct {
	*ControllerCommon
}

func NewControllerServer(common *ControllerCommon,
	stream *stream.StreamService) *ControllerServer {
	server := &ControllerServer{
		ControllerCommon: common,
	}

	// server
	stream.Subscribe("register_server", server.RegisterServer)

	return server
}

func (c *ControllerServer) RegisterServer(client *stream.Client, value interface{}) {

	log.Info("call register server")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	token, err := c.endpoint.RegisterServer()
	if err != nil {
		c.Err(client, value, err)
		return
	}

	client.Token = token

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "token": token})
	client.Send <- msg

	return
}

func (c *ControllerServer) RemoveServerToken(client *stream.Client, value interface{}) {

	log.Info("call remove server")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, value, err)
		return
	}

	if err = c.endpoint.RemoveServerToken(server); err != nil {
		c.Err(client, value, err)
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send <- msg
}
