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

package controllers

import (
	m "github.com/e154/smart-home-gate/models"
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

	var accessToken string
	if client.Id == "" {
		server, err := c.endpoint.RegisterServer()
		if err != nil {
			c.Err(client, message, err)
			return
		}
		accessToken = server.GenAccessToken()

		client.Id = server.Id
		client.Token = server.Token

		log.Infof("register new server: %s", server.Id)

	} else {
		server := &m.Server{
			Id:    client.Id,
			Token: client.Token,
		}
		accessToken = server.GenAccessToken()
		log.Infof("use client id: %s", client.Id)
	}

	log.Infof("accessToken %s", accessToken)

	payload := map[string]interface{}{
		"token": accessToken,
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
