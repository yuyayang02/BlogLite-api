package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
)

func BindJSON[T any](c *gin.Context) (T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, errors.InvalidParamsError(err)
	}
	return req, nil
}
