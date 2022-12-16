// Code generated by 'yaegi extract github.com/kuchensheng/bintools/json/model'. DO NOT EDIT.

package service

import (
	"github.com/kuchensheng/bintools/json/model"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["github.com/kuchensheng/bintools/json/model/model"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"GenerateFile2Go":              reflect.ValueOf(model.GenerateFile2Go),
		"GenerateJson2Go":              reflect.ValueOf(model.GenerateJson2Go),
		"NewBusinessException":         reflect.ValueOf(model.NewBusinessException),
		"NewBusinessExceptionWithData": reflect.ValueOf(model.NewBusinessExceptionWithData),
		"PLUGIN_PATH":                  reflect.ValueOf(constant.MakeFromLiteral("\"plugins\"", token.STRING, 0)),

		// type definitions
		"ApiStepPredicate":    reflect.ValueOf((*model.ApiStepPredicate)(nil)),
		"ApixApi":             reflect.ValueOf((*model.ApixApi)(nil)),
		"ApixData":            reflect.ValueOf((*model.ApixData)(nil)),
		"ApixParameter":       reflect.ValueOf((*model.ApixParameter)(nil)),
		"ApixProperty":        reflect.ValueOf((*model.ApixProperty)(nil)),
		"ApixResponse":        reflect.ValueOf((*model.ApixResponse)(nil)),
		"ApixRule":            reflect.ValueOf((*model.ApixRule)(nil)),
		"ApixSchema":          reflect.ValueOf((*model.ApixSchema)(nil)),
		"ApixScript":          reflect.ValueOf((*model.ApixScript)(nil)),
		"ApixSetCookie":       reflect.ValueOf((*model.ApixSetCookie)(nil)),
		"ApixStep":            reflect.ValueOf((*model.ApixStep)(nil)),
		"ApixSwitchPredicate": reflect.ValueOf((*model.ApixSwitchPredicate)(nil)),
		"BusinessException":   reflect.ValueOf((*model.BusinessException)(nil)),
	}
}