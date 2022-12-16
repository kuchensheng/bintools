//go:generate yaegi extract "github.com/gin-gonic/gin"
//go:generate yaegi extract "github.com/kuchensheng/bintools/tracer/trace"
//go:generate yaegi extract "github.com/rs/zerolog/log"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/model"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/consts"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/executor/js"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/executor/parameter"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/executor/predicate"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/executor/server"
//go:generate yaegi extract "github.com/kuchensheng/bintools/json/executor/response"

package service

import "reflect"

// Symbols variable stores the map of stdlib symbols per package.
var Symbols = map[string]map[string]reflect.Value{}

func init() {
	Symbols["github.com/traefik/yaegi/stdlib/stdlib"] = map[string]reflect.Value{
		"Symbols": reflect.ValueOf(Symbols),
	}
}
