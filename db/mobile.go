package db

import (
	"github.com/e154/smart-home-gate/system/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Mobiles struct {
	Db *gorm.DB
}

type Mobile struct {
	Id        int64 `gorm:"primary_key"`
	Token     uuid.UUID
	ServerId  int64
	Server    *Server
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Mobile) TableName() string {
	return "mobiles"
}

func (n Mobiles) Add(mobile *Mobile) (id int64, err error) {

	if err = n.Db.Create(mobile).Error; err != nil {
		return
	}
	id = mobile.Id
	return
}

func (n *Mobiles) Update(m *Mobile) (err error) {
	err = n.Db.Model(&Mobile{}).
		Updates(map[string]interface{}{
			"token": m.Token,
		}).
		Where("id = ?", m.Id).
		Error
	return
}

func (n *Mobiles) Remove(mobileId int64) (err error) {
	err = n.Db.Delete(&Mobile{Id: mobileId}).Error
	return
}

func (n Mobiles) GetById(id int64) (mobile *Mobile, err error) {
	mobile = &Mobile{Id: id}
	err = n.Db.First(&mobile).Error
	return
}

func (n Mobiles) GetByToken(token string) (mobile *Mobile, err error) {
	mobile = &Mobile{}
	err = n.Db.Model(mobile).
		Where("token = ?", token).
		First(&mobile).
		Error
	return
}

func (n *Mobiles) List(limit, offset int64) (list []*Mobile, total int64, err error) {

	if err = n.Db.Model(Mobile{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Mobile, 0)

	err = n.Db.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	return
}
