package models

import (
	"github.com/e154/smart-home-gate/system/uuid"
	"net"
	"time"
)

type Client struct {
	Id               int64     `json:"id"`
	ClientId         uuid.UUID `json:"client_id"`
	Token            string    `json:"token"`
	Ip               net.IP    `json:"ip"`
	TokenGeneratedAt time.Time `json:"token_generated_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (c *Client) GenToken() *Client {

	return c
}
