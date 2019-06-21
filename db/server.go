package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Servers struct {
	Db *gorm.DB
}

type Server struct {
	Id        int64 `gorm:"primary_key"`
	Token     string
	Mobiles   []*Mobile
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Server) TableName() string {
	return "servers"
}

func (n Servers) Add(server *Server) (id int64, err error) {
	if err = n.Db.Create(&server).Error; err != nil {
		return
	}
	id = server.Id
	return
}

func (n *Servers) Update(m *Server) (err error) {
	err = n.Db.Model(&Server{}).
		Updates(map[string]interface{}{
			"token": m.Token,
		}).
		Where("id = ?", m.Id).
		Error
	return
}

func (n *Servers) Remove(serverId int64) (err error) {
	err = n.Db.Delete(&Server{Id: serverId}).Error
	return
}

func (n Servers) GetById(id int64) (server *Server, err error) {
	server = &Server{Id: id}
	if err = n.Db.First(&server).Error; err != nil {
		return
	}
	err = n.DependencyLoading(server)
	return
}

func (n Servers) GetByToken(token string) (server *Server, err error) {
	server = &Server{}
	err = n.Db.Model(server).
		Where("token = ?", token).
		First(&server).
		Error
	if err != nil {
		return
	}
	err = n.DependencyLoading(server)
	return
}

func (n *Servers) List(limit, offset int64) (list []*Server, total int64, err error) {

	if err = n.Db.Model(Server{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Server, 0)

	err = n.Db.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	if err != nil {
		return
	}

	for _, s := range list {
		err = n.DependencyLoading(s)
	}
	return
}

func (n *Servers) DependencyLoading(server *Server) (err error) {
	server.Mobiles = make([]*Mobile, 0)
	n.Db.Model(server).
		Related(&server.Mobiles)
	return
}
