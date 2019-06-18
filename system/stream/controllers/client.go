package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/system/stream"
	"reflect"
)

type ControllerClient struct {
	adaptors *adaptors.Adaptors
}

func NewControllerClient(adaptors *adaptors.Adaptors) *ControllerClient {
	return &ControllerClient{
		adaptors: adaptors,
	}
}

func (c *ControllerClient) GetToken(client *stream.Client, value interface{}) {

	fmt.Print("GetToken")

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	//var filter string
	//if filter, ok = v["filter"].(string); ok {
	//}

	//images, err := c.adaptors.Image.GetAllByDate(filter)
	//if err != nil {
	//	client.Notify("error", err.Error())
	//	log.Error(err.Error())
	//	return
	//}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "token": "xxxxxxxxxx"})
	client.Send <- msg

	return
}
