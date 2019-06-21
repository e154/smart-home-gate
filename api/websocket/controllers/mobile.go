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

	// mobile
	stream.Subscribe("register_mobile", mobile.RegisterMobile)

	return mobile
}

func (c *ControllerMobile) RegisterMobile(client *stream.Client, value interface{}) {

	log.Info("call register mobile")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	token, err := c.endpoint.RegisterMobile()
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
