package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/yuyayang02/BlogLite-api/internal/common/config"
	"github.com/yuyayang02/BlogLite-api/internal/common/constant"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"github.com/yuyayang02/BlogLite-api/internal/common/server/httpresponse"
	"github.com/yuyayang02/BlogLite-api/internal/common/utils"
)

type GetAuthJSONParams struct {
	Password string `json:"password"`
}

func GetAuth(c *gin.Context) {
	req, err := utils.BindJSON[GetAuthJSONParams](c)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	if req.Password != config.Cfg.AuthAdminPassword {
		httpresponse.Error(c, errors.PWDError)
		return
	}

	sign, err := Sign(NewUserCliaims(Admin.ID, Admin.Type, Admin.Name, constant.DefaultJWTAuthDuration))
	if err != nil {
		httpresponse.Error(c, errors.InternalServiceError(err.Error()))
		return
	}

	c.Header("X-Auth-Token", sign)
	httpresponse.OK(c)
}
