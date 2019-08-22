package endpoint

import (
	"errors"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/uuid"
)

type Mobile struct {
	*CommonEndpoint
}

func NewMobile(common *CommonEndpoint) *Mobile {
	return &Mobile{CommonEndpoint: common}
}

func (c *Mobile) RegisterMobile(server *m.Server) (token string, err error) {

	mobileClient := &m.Mobile{
		Token:    uuid.NewV4(),
		ServerId: server.Id,
	}

	if _, err = c.adaptors.Mobile.Add(mobileClient); err != nil {
		return
	}

	token = mobileClient.GenAccessToken()

	return
}

func (c *Mobile) RemoveMobileToken(server *m.Server, token string) (err error) {

	mobile, err := c.adaptors.Mobile.GetByAccessToken(token)
	if err != nil {
		return
	}

	var exist bool
	for _, m := range server.Mobiles {
		if m.Token == mobile.Token {
			exist = true
		}
	}

	if !exist {
		err = errors.New("mobile not found")
		return
	}


	err = c.adaptors.Mobile.Remove(mobile)

	return
}

func (c *Mobile) ListMobileToken(limit, offset int64) (list []*m.Mobile, total int64, err error) {
	list, total, err = c.adaptors.Mobile.List(limit, offset)
	return
}
