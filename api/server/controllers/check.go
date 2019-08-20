package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerCheck struct {
	*ControllerCommon
}

func NewControllerCheck(common *ControllerCommon) *ControllerCheck {
	return &ControllerCheck{ControllerCommon: common}
}

// swagger:operation GET /check/mobile_access_token check
// ---
// summary: mobile access token check page
// description:
// security:
// - ServerAuthorization: []
// consumes:
// - text/plain
// produces:
// - text/plain
// tags:
// - check
// responses:
//   "200":
//	   description: Success response
//   "400":
//	   description: Bad request
//   "404":
//	   description: Not found
//
func (i ControllerCheck) CheckMobileAccessToken(ctx *gin.Context) {

	accessToken := ctx.GetHeader("ServerAuthorization")
	if accessToken == "" {
		ctx.String(http.StatusBadRequest, "need access token")
		return
	}

	ok := i.endpoint.CheckMobileAccessToken(accessToken)
	fmt.Println("ok", ok)
	if ok {
		ctx.String(http.StatusOK, "ok")
	} else {
		ctx.String(http.StatusNotFound, "not found")
	}

	return
}

// swagger:operation GET /check/mobile_access check
// ---
// summary: mobile access connection page
// description:
// consumes:
// - text/plain
// produces:
// - text/plain
// tags:
// - check
// responses:
//   "200":
//	   description: Success response
//
func (i ControllerCheck) CheckMobileAccess(ctx *gin.Context) {
	apiVersion := "smart-home-gate"
	ctx.String(http.StatusOK, apiVersion)
	return
}
