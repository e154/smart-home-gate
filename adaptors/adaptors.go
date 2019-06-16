package adaptors

import (
	"github.com/e154/smart-home-gate/system/config"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("adaptors")
)

type Adaptors struct {
}

func NewAdaptors(db *gorm.DB,
	cfg *config.AppConfig) (adaptors *Adaptors) {

	adaptors = &Adaptors{

	}

	return
}
