package models

import (
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"github.com/e154/smart-home-gate/system/uuid"
	"time"
)

type Mobile struct {
	Id        string    `json:"id"`
	Token     uuid.UUID `json:"token"`
	ServerId  string    `json:"server_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Mobile) GenAccessToken() (token string) {

	var requestId = common.RandomString(10, common.Charset)
	var timestamp = s.CreatedAt.Unix()
	token = fmt.Sprintf("%s-%s-%d-%s", s.Id, requestId, timestamp, common.Sha256(requestId+s.Token.String()+fmt.Sprintf("%d", timestamp)))

	return
}
