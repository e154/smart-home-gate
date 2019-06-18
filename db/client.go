package db

import (
	"github.com/e154/smart-home/system/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Clients struct {
	Db *gorm.DB
}

type Client struct {
	Id        int64 `gorm:"primary_key"`
	ClientId  uuid.UUID
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
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

func (n Clients) GetById(clientId int64) (client *Client, err error) {
	client = &Client{}
	err = n.Db.Model(client).
		Where("client_id = ?", clientId).
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
