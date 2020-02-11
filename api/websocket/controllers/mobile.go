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
	"github.com/e154/smart-home-gate/system/stream"
)

type ControllerMobile struct {
	*ControllerCommon
}

func NewControllerMobile(common *ControllerCommon,
	stream *stream.StreamService) *ControllerMobile {
	mobile := &ControllerMobile{
		ControllerCommon: common,
	}

	// register methods
	stream.Subscribe("register_mobile", mobile.RegisterMobile)
	stream.Subscribe("remove_mobile", mobile.RemoveMobileToken)
	stream.Subscribe("mobile_token_list", mobile.ListMobileTokens)

	return mobile
}

func (c *ControllerMobile) RegisterMobile(client *stream.Client, message stream.Message) {

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	log.Info("call register mobile")

	token, err := c.endpoint.RegisterMobile(server)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	payload := map[string]interface{}{
		"token": token,
	}
	response := message.Response(payload)

	client.Send <- response.Pack()

	return
}

func (c *ControllerMobile) RemoveMobileToken(client *stream.Client, message stream.Message) {

	log.Info("call remove mobile")

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	token := message.Payload["token"].(string)

	if err = c.endpoint.RemoveMobileToken(server, token); err != nil {
		c.Err(client, message, err)
		return
	}

	client.Send <- message.Success().Pack()
}

func (c *ControllerMobile) ListMobileTokens(client *stream.Client, message stream.Message) {

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	tokenList := make([]string, 0)
	for _, m := range server.Mobiles {
		tokenList = append(tokenList, m.GenAccessToken())
	}

	payload := map[string]interface{}{
		"total":      len(tokenList),
		"token_list": tokenList,
	}
	response := message.Response(payload)

	client.Send <- response.Pack()
}
