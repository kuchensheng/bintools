package predicate

import (
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/model"
)

//ExecPredicates predicateType = 1 所有条件都为真，否则任一条件为真
func ExecPredicates(ctx *gin.Context, predicates []model.ApiStepPredicate, predicateType int) (bool, error) {
	return true, nil
}
