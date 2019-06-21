package adaptors

import (
	"github.com/e154/smart-home-gate/system/config"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

var (
	log      = logging.MustGetLogger("adaptors")
	HashSalt string
)

type Adaptors struct {
	Server   *Server
	Mobile   *Mobile
	Variable *Variable
}

func NewAdaptors(db *gorm.DB,
	cfg *config.AppConfig) (adaptors *Adaptors) {

	adaptors = &Adaptors{
		Server:   GetServerAdaptor(db),
		Mobile:   GetMobileAdaptor(db),
		Variable: GetVariableAdaptor(db),
	}

	return
}
