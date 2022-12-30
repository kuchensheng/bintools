package check

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
)

func CheckTenantId(ctx *gin.Context, tenantId string) error {
	if tId, ok := ctx.Get(consts.TENANT_ID); !ok {
		return errors.New("禁止操作：租户不匹配")
	} else if tenantId != tId.(string) {
		return errors.New("禁止操作：租户不匹配")
	}
	return nil
}
