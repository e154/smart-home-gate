package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/e154/smart-home-gate/system/config"
	"github.com/op/go-logging"
	"github.com/speps/go-hashids"
	"math/rand"
	"time"
)

var (
	log        = logging.MustGetLogger("common")
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandomString(length int, charset string) string {
	b := make([]byte, length*2)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetHashFromId(id int64, salt ...interface{}) (hash string, err error) {

	hd := hashids.NewData()
	hd.Salt = config.HashSalt
	hd.MinLength = 4
	switch len(salt) {
	case 1:
		hd.Salt = salt[0].(string)
	case 2:
		hd.Salt = salt[0].(string)
		hd.MinLength = salt[1].(int)
	}

	h, _ := hashids.NewWithData(hd)
	hash, err = h.EncodeInt64([]int64{id})

	return
}

func GetIdFromHash(hash string, salt ...interface{}) (id int64, err error) {

	hd := hashids.NewData()
	hd.Salt = config.HashSalt
	hd.MinLength = 4
	switch len(salt) {
	case 1:
		hd.Salt = salt[0].(string)
	case 2:
		hd.Salt = salt[0].(string)
		hd.MinLength = salt[1].(int)
	}

	h, _ := hashids.NewWithData(hd)

	var ids []int64
	if ids, err = h.DecodeInt64WithError(hash); err != nil {
		return
	}

	if len(ids) > 0 {
		id = ids[0]
	}

	return
}

func Sha256(input string) string {
	sha_256 := sha256.New()
	sha_256.Write([]byte(input))
	return hex.EncodeToString(sha_256.Sum(nil))
}

func ComputeHmac256() string {
	var message = "token"
	var secret = RandomString(255, Charset)

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}