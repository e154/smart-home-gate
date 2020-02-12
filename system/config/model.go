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

type AppConfig struct {
	ServerHost        string  `json:"server_host"`
	ServerPort        int     `json:"server_port"`
	AutoMigrate       bool    `json:"auto_migrate"`
	PgUser            string  `json:"pg_user"`
	PgPass            string  `json:"pg_pass"`
	PgHost            string  `json:"pg_host"`
	PgName            string  `json:"pg_name"`
	PgPort            string  `json:"pg_port"`
	PgDebug           bool    `json:"pg_debug"`
	PgLogger          bool    `json:"pg_logger"`
	PgMaxIdleConns    int     `json:"pg_max_idle_conns"`
	PgMaxOpenConns    int     `json:"pg_max_open_conns"`
	PgConnMaxLifeTime int     `json:"pg_conn_max_life_time"`
	HashSalt          string  `json:"hash_salt"`
	MetricPort        int     `json:"metric_port"`
	Mode              RunMode `json:"mode"`
}

type RunMode string

const (
	DebugMode   = RunMode("debug")
	ReleaseMode = RunMode("release")
)
