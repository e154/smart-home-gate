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
	"github.com/e154/smart-home-gate/system/config"
	"github.com/e154/smart-home-gate/system/migrations"
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
	cfg *config.AppConfig,
	migrations *migrations.Migrations) (adaptors *Adaptors) {

	if cfg.AutoMigrate {
		migrations.Up()
	}

	adaptors = &Adaptors{
		Server:   GetServerAdaptor(db),
		Mobile:   GetMobileAdaptor(db),
		Variable: GetVariableAdaptor(db),
	}

	return
}
