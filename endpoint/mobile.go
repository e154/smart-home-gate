package endpoint

import (
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

func (c *Mobile) RemoveMobileToken(mobile *m.Mobile) (err error) {
	err = c.adaptors.Mobile.Remove(mobile)
	return
}

func (c *Mobile) ListMobileToken(limit, offset int64) (list []*m.Mobile, total int64, err error) {
	list, total, err = c.adaptors.Mobile.List(limit, offset)
	return
}
