package endpoint

import (
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"strconv"
	"strings"
)

type Check struct {
	*CommonEndpoint
}

func NewCheck(common *CommonEndpoint) *Check {
	return &Check{CommonEndpoint: common}
}

func (c *Check) CheckMobileAccessToken(accessToken string) (ok bool) {

	data := strings.Split(accessToken, "-")
	if len(data) != 4 {
		return
	}

	mobileClientId := data[0]
	requestRandomId := data[1]
	hash := data[3]

	timestamp, errw := strconv.Atoi(data[2])
	if errw != nil {
		return
	}

	if len(requestRandomId) < 8 {
		return
	}

	mobileObj, err := c.adaptors.Mobile.GetById(mobileClientId)
	if err != nil {
		return
	}

	if hash != common.Sha256(requestRandomId+mobileObj.Token.String()+fmt.Sprintf("%d", timestamp)) {
		return
	}

	ok = true

	return
}
