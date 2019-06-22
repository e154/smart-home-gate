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

	client.Token = token

	payload := map[string]interface{}{
		"token": token,
	}
	response := message.Response(payload)

	client.Send <- response.Pack()

	return
}

func (c *ControllerMobile) RemoveMobileToken(client *stream.Client, message stream.Message) {

	log.Info("call remove mobile")

	mobile, err := c.GetMobile(client)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	if err = c.endpoint.RemoveMobileToken(mobile); err != nil {
		c.Err(client, message, err)
		return
	}

	client.Send <- message.Success().Pack()
}

func (c *ControllerMobile) ListMobileTokens(client *stream.Client, message stream.Message) {

	list, total, err := c.endpoint.ListMobileToken(99, 0)
	if err != nil {
		c.Err(client, message, err)
		return
	}

	tokenList := make([]string, 0)
	for _, m := range list {
		tokenList = append(tokenList, m.Token.String())
	}

	payload := map[string]interface{}{
		"total":      total,
		"token_list": tokenList,
	}
	response := message.Response(payload)

	client.Send <- response.Pack()
}
