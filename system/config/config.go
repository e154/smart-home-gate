package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const path = "conf/config.json"

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
}
