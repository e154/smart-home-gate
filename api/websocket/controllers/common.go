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
	"fmt"
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/endpoint"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/stream"
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	stream   *stream.StreamService
	endpoint endpoint.IEndpoint
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint endpoint.IEndpoint) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		endpoint: endpoint,
		stream:   stream,
	}
}

func (c *ControllerCommon) GetServer(client *stream.Client) (server *m.Server, err error) {

	if client == nil {
		err = fmt.Errorf("nil client")
		return
	}

	if client.Token == "" {
		err = fmt.Errorf("zero server token")
		return
	}

	server, err = c.adaptors.Server.GetById(client.Id)

	return
}

func (c *ControllerCommon) GetMobile(client *stream.Client) (mobile *m.Mobile, err error) {

	if client == nil {
		err = fmt.Errorf("nil client")
		return
	}

	if client.Token == "" {
		err = fmt.Errorf("zero mobile token")
		return
	}

	mobile, err = c.adaptors.Mobile.GetByAccessToken(client.Token)

	return
}

func (c *ControllerCommon) Err(client *stream.Client, message stream.Message, err error) {
	msg := stream.Message{
		Id: message.Id,
		Forward: stream.Response,
		Status: stream.StatusError,
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	}
	client.Send <- msg.Pack()
}
