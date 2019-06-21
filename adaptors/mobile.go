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

	dbVer := n.toDb(ver)
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
	}
	if _, err = n.table.GetById(id); err != nil {
		return
	}

	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Mobile) Remove(ver *m.Mobile) (err error) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}
	err = n.table.Remove(id)
	return
}

func (n *Mobile) GetById(verId string) (ver *m.Mobile, err error) {

	id, err := common.GetIdFromHash(verId, HashSalt)
	if err != nil {
		log.Error(err.Error())
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

	if hash != common.Sha256(requestRandomId+dbVer.Token+fmt.Sprintf("%d", timestamp)) {
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
	ver = &m.Mobile{
		Id:        id,
		Token:     dbVer.Token,
		ServerId:  dbVer.ServerId,
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	return
}

func (n *Mobile) toDb(ver *m.Mobile) (dbVer *db.Mobile) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}
	dbVer = &db.Mobile{
		Id:        id,
		Token:     ver.Token,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
	return
}
