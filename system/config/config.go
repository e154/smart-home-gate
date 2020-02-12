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

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const path = "conf/config.json"

var (
	HashSalt string
)

func ReadConfig() (conf *AppConfig, err error) {
	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}
	conf = &AppConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
		return
	}

	checkEnv(conf)

	//TODO fix
	HashSalt = conf.HashSalt

	return
}

func checkEnv(conf *AppConfig) {

	if serverHost := os.Getenv("SERVER_HOST"); serverHost != "" {
		conf.ServerHost = serverHost
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		v, _ := strconv.ParseInt(serverPort, 10, 32)
		conf.ServerPort = int(v)
	}

	if pgUser := os.Getenv("PG_USER"); pgUser != "" {
		conf.PgUser = pgUser
	}

	if pgPass := os.Getenv("PG_PASS"); pgPass != "" {
		conf.PgPass = pgPass
	}

	if pgHost := os.Getenv("PG_HOST"); pgHost != "" {
		conf.PgHost = pgHost
	}

	if pgName := os.Getenv("PG_NAME"); pgName != "" {
		conf.PgName = pgName
	}

	if pgPort := os.Getenv("PG_PORT"); pgPort != "" {
		conf.PgPort = pgPort
	}

	if hashSalt := os.Getenv("HashSalt"); hashSalt != "" {
		conf.HashSalt = hashSalt
	}

	if pgDebug := os.Getenv("PG_DEBUG"); pgDebug != "" {
		conf.PgDebug, _ = strconv.ParseBool(pgDebug)
	}

	if pgLogger := os.Getenv("PG_LOGGER"); pgLogger != "" {
		conf.PgLogger, _ = strconv.ParseBool(pgLogger)
	}

	if pgMaxIdleConns := os.Getenv("PG_MAX_IDLE_CONNS"); pgMaxIdleConns != "" {
		v, _ := strconv.ParseInt(pgMaxIdleConns, 10, 32)
		conf.PgMaxIdleConns = int(v)
	}

	if pgMaxOpenConns := os.Getenv("PG_MAX_OPEN_CONNS"); pgMaxOpenConns != "" {
		v, _ := strconv.ParseInt(pgMaxOpenConns, 10, 32)
		conf.PgMaxOpenConns = int(v)
	}

	if pgConnMaxLifeTime := os.Getenv("PG_CONN_MAX_LIFE_TIME"); pgConnMaxLifeTime != "" {
		v, _ := strconv.ParseInt(pgConnMaxLifeTime, 10, 32)
		conf.PgConnMaxLifeTime = int(v)
	}

	if autoMigrate := os.Getenv("AUTO_MIGRATE"); autoMigrate != "" {
		conf.AutoMigrate, _ = strconv.ParseBool(autoMigrate)
	}

	if metricPort := os.Getenv("METRIC_PORT"); metricPort != "" {
		v, _ := strconv.ParseInt(metricPort, 10, 32)
		conf.MetricPort = int(v)
	}

	if mode := os.Getenv("MODE"); mode != "" {
		conf.Mode = RunMode(mode)
	}
}
