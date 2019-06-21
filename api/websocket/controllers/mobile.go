package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home-gate/system/stream"
	"reflect"
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

func (c *ControllerMobile) RegisterMobile(client *stream.Client, value interface{}) {

	server, err := c.GetServer(client)
	if err != nil {
		c.Err(client, value, err)
		return
	}

	log.Info("call register mobile")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	token, err := c.endpoint.RegisterMobile(server)
	if err != nil {
		c.Err(client, value, err)
		return
	}

	client.Token = token

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "token": token})
	client.Send <- msg

	return
}

func (c *ControllerMobile) RemoveMobileToken(client *stream.Client, value interface{}) {

	log.Info("call remove mobile")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	mobile, err := c.GetMobile(client)
	if err != nil {
		c.Err(client, value, err)
		return
	}

	if err = c.endpoint.RemoveMobileToken(mobile); err != nil {
		c.Err(client, value, err)
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send <- msg
}

func (c *ControllerMobile) ListMobileTokens(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	list, total, err := c.endpoint.ListMobileToken(99, 0)
	if err != nil {
		c.Err(client, value, err)
		return
	}

	tokenList := make([]string, 0)
	for _, m := range list {
		tokenList = append(tokenList, m.Token.String())
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok", "total": total, "token_list": tokenList})
	client.Send <- msg
}