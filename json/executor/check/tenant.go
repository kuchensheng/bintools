package check

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
)

func CheckTenantId(ctx *gin.Context, tenantId string) error {
	tId := ctx.GetHeader(consts.TENANT_ID)
	if tenantId != tId {
		return errors.New("禁止操作：租户不匹配")
	}
	return nil
}
