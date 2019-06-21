package models

import (
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"time"
)

type Server struct {
	Id        string    `json:"id"`
	Token     string    `json:"token"`
	Mobiles   []*Mobile `json:"mobiles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Server) GenAccessToken() (token string) {

	var requestId = common.RandomString(10, common.Charset)
	var timestamp = time.Now().Unix()
	token = fmt.Sprintf("%s-%s-%d-%s", s.Id, requestId, timestamp, common.Sha256(requestId+s.Token+fmt.Sprintf("%d", timestamp)))

	return
}
