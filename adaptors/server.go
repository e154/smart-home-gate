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

type Server struct {
	table *db.Servers
	db    *gorm.DB
}

func GetServerAdaptor(d *gorm.DB) *Server {
	return &Server{
		table: &db.Servers{Db: d},
		db:    d,
	}
}

func (n *Server) Add(ver *m.Server) (idStr string, err error) {

	dbVer := n.toDb(ver)
	var id int64
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	idStr, err = common.GetHashFromId(id, HashSalt)
	ver.Id = idStr

	return
}

func (n *Server) Update(ver *m.Server) (err error) {
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

func (n *Server) Remove(ver *m.Server) (err error) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}
	err = n.table.Remove(id)
	return
}

func (n *Server) GetById(verId string) (ver *m.Server, err error) {

	id, err := common.GetIdFromHash(verId, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}

	var dbVer *db.Server
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Server) GetByToken(token string) (ver *m.Server, err error) {

	var dbVer *db.Server
	if dbVer, err = n.table.GetByToken(token); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Server) GetByAccessToken(accessToken string) (ver *m.Server, err error) {

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

	var dbVer *db.Server
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

func (n *Server) List(limit, offset int64) (list []*m.Server, total int64, err error) {

	var dbList []*db.Server
	if dbList, total, err = n.table.List(limit, offset); err != nil {
		return
	}

	list = make([]*m.Server, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Server) fromDb(dbVer *db.Server) (ver *m.Server) {
	id, err := common.GetHashFromId(dbVer.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}
	ver = &m.Server{
		Id:        id,
		Token:     dbVer.Token,
		Mobiles:   make([]*m.Mobile, 0),
		CreatedAt: dbVer.CreatedAt,
		UpdatedAt: dbVer.UpdatedAt,
	}

	// Mobiles
	if len(dbVer.Mobiles) > 0 {
		mobileAdaptor := GetMobileAdaptor(n.db)
		for _, dbConn := range dbVer.Mobiles {
			con := mobileAdaptor.fromDb(dbConn)
			ver.Mobiles = append(ver.Mobiles, con)
		}
	}

	return
}

func (n *Server) toDb(ver *m.Server) (dbVer *db.Server) {
	id, err := common.GetIdFromHash(ver.Id, HashSalt)
	if err != nil {
		log.Error(err.Error())
	}
	dbVer = &db.Server{
		Id:        id,
		Token:     ver.Token,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
	return
}
