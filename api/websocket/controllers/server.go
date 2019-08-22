package controllers

import (
	"github.com/e154/smart-home-gate/system/stream"
)

type ControllerServer struct {
	*ControllerCommon
}

func NewControllerServer(common *ControllerCommon,
	stream *stream.StreamService) *ControllerServer {
	server := &ControllerServer{
		ControllerCommon: common,
	}

	// register methods
	stream.Subscribe("register_server", server.RegisterServer)
	stream.Subscribe("remove_server", server.RemoveServerToken)

	return server
}

func (c *ControllerServer) RegisterServer(client *stream.Client, message stream.Message) {

	log.Info("call register server")

	token, err := c.endpoint.RegisterServer()
	if err != nil {
		c.Err(client, message, err)
		return
	}

	client.Token = token

	payload := map[string]interface{}{
		"token": token,
	}
	response := message.Response(payload)

	client.Send <- response.Pack()

	return
}

func (c *ControllerServer) RemoveServerToken(client *stream.Client, message stream.Message) {

	log.Info("call remove server")

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	if err = c.endpoint.RemoveServerToken(server); err != nil {
		c.Err(client, message, err)
		return
	}

	client.Send <- message.Success().Pack()
}
