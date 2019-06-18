package adaptors

import (
	"github.com/e154/smart-home-gate/db"
	m "github.com/e154/smart-home-gate/models"
	"github.com/jinzhu/gorm"
)

type Client struct {
	table *db.Clients
	db    *gorm.DB
}

func GetClientAdaptor(d *gorm.DB) *Client {
	return &Client{
		table: &db.Clients{Db: d},
		db:    d,
	}
}

func (n *Client) Add(ver *m.Client) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *Client) Update(ver *m.Client) (err error) {

	if _, err = n.table.GetById(ver.Id); err != nil {
		return
	}

	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Client) GetById(verId int64) (ver *m.Client, err error) {

	var dbVer *db.Client
	if dbVer, err = n.table.GetById(verId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Client) GetByToken(token string) (ver *m.Client, err error) {

	var dbVer *db.Client
	if dbVer, err = n.table.GetByToken(token); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Client) List(limit, offset int64) (list []*m.Client, total int64, err error) {

	var dbList []*db.Client
	if dbList, total, err = n.table.List(limit, offset); err != nil {
		return
	}

	list = make([]*m.Client, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Client) fromDb(dbVer *db.Client) (ver *m.Client) {
	ver = &m.Client{
		Id:               dbVer.Id,
		ClientId:         dbVer.ClientId,
		Token:            dbVer.Token,
		Ip:               dbVer.Ip,
		TokenGeneratedAt: dbVer.TokenGeneratedAt,
		CreatedAt:        dbVer.CreatedAt,
		UpdatedAt:        dbVer.UpdatedAt,
	}

	return
}

func (n *Client) toDb(ver *m.Client) (dbVer *db.Client) {
	dbVer = &db.Client{
		Id:               ver.Id,
		ClientId:         ver.ClientId,
		Token:            ver.Token,
		Ip:               ver.Ip,
		TokenGeneratedAt: ver.TokenGeneratedAt,
		CreatedAt:        ver.CreatedAt,
		UpdatedAt:        ver.UpdatedAt,
	}
	return
}
