package db

import (
	"github.com/e154/smart-home-gate/system/uuid"
	"github.com/jinzhu/gorm"
	"net"
	"time"
)

type Clients struct {
	Db *gorm.DB
}

type Client struct {
	Id               int64 `gorm:"primary_key"`
	ClientId         uuid.UUID
	Token            string
	TokenGeneratedAt time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Ip               net.IP
}

func (m *Client) TableName() string {
	return "clients"
}

func (n Clients) Add(client *Client) (id int64, err error) {
	if err = n.Db.Create(&client).Error; err != nil {
		return
	}
	id = client.Id
	return
}

func (n *Clients) Update(m *Client) (err error) {
	err = n.Db.Model(&Client{}).
		Updates(map[string]interface{}{
			"token":              m.Token,
			"token_generated_at": m.TokenGeneratedAt,
		}).
		Where("client_id = ?", m.ClientId).
		Error
	return
}

func (n Clients) GetById(id int64) (client *Client, err error) {
	client = &Client{Id: id}
	err = n.Db.First(&client).Error
	return
}

func (n Clients) GetByToken(token string) (client *Client, err error) {
	client = &Client{}
	err = n.Db.Model(client).
		Where("token = ?", token).
		First(&client).
		Error
	return
}

func (n *Clients) List(limit, offset int64) (list []*Client, total int64, err error) {

	if err = n.Db.Model(Client{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Client, 0)

	err = n.Db.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	return
}
