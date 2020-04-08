// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package adaptors

import (
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"github.com/e154/smart-home-gate/db"
	m "github.com/e154/smart-home-gate/models"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type Mobile struct {
	table *db.Mobiles
	db    *gorm.DB
}

func GetMobileAdaptor(d *gorm.DB) *Mobile {
	return &Mobile{
		table: &db.Mobiles{Db: d},
		db:    d,
	}
}

func (n *Mobile) Add(ver *m.Mobile) (idStr string, err error) {

	var dbVer *db.Mobile
	if dbVer, err = n.toDb(ver); err != nil {
		return
	}

	var id int64
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	idStr, err = common.GetHashFromId(id, HashSalt)
	ver.Id = idStr

	return
}

func (n *Mobile) Update(ver *m.Mobile) (err error) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if _, err = n.table.GetById(id); err != nil {
		return
	}

	var dbVer *db.Mobile
	dbVer, err = n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Mobile) Remove(ver *m.Mobile) (err error) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = n.table.Remove(id)
	return
}

func (n *Mobile) GetById(verId string) (ver *m.Mobile, err error) {

	id, err := common.GetIdFromHash(verId, HashSalt)
	if err != nil {
		log.Error(err.Error())
		return
	}

	var dbVer *db.Mobile
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Mobile) GetByToken(token string) (ver *m.Mobile, err error) {

	var dbVer *db.Mobile
	if dbVer, err = n.table.GetByToken(token); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Mobile) GetByAccessToken(accessToken string) (ver *m.Mobile, err error) {

	data := strings.Split(accessToken, "-")
	if len(data) != 4 {
		err = fmt.Errorf("access token not valid")
		return
	}

	id, err := common.GetIdFromHash(data[0], HashSalt)
	if err != nil {
		log.Error(err.Error())
		return
	}

	requestRandomId := data[1]
	hash := data[3]

	timestamp, errw := strconv.Atoi(data[2])
	if errw != nil {
		err = fmt.Errorf("access token not valid")
		return
	}

	if len(requestRandomId) < 8 {
		err = fmt.Errorf("access token not valid")
		return
	}

	var dbVer *db.Mobile
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	if hash != common.Sha256(requestRandomId+dbVer.Token.String()+fmt.Sprintf("%d", timestamp)) {
		err = fmt.Errorf("Wrong auth data, wrong hash")
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Mobile) List(limit, offset int64) (list []*m.Mobile, total int64, err error) {

	var dbList []*db.Mobile
	if dbList, total, err = n.table.List(limit, offset); err != nil {
		return
	}

	list = make([]*m.Mobile, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Mobile) fromDb(dbVer *db.Mobile) (ver *m.Mobile) {
	id, err := common.GetHashFromId(dbVer.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}

	serverId, err := common.GetHashFromId(dbVer.ServerId, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}

	ver = &m.Mobile{
		Id:        id,
		Token:     dbVer.Token,
		ServerId:  serverId,
		RequestId: dbVer.RequestId,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	return
}

func (n *Mobile) toDb(ver *m.Mobile) (dbVer *db.Mobile, err error) {
	var id int64
	if id, err = common.GetIdFromHash(ver.Id, HashSalt); err != nil {
		log.Error(err.Error())
		return
	}

	var serverId int64
	if serverId, err = common.GetIdFromHash(ver.ServerId, HashSalt); err != nil {
		log.Error(err.Error())
		return
	}
	dbVer = &db.Mobile{
		Id:        id,
		ServerId:  serverId,
		Token:     ver.Token,
		RequestId: ver.RequestId,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}

	return
}
