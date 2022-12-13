//go:generate yaegi extract github.com/gin-gonic/gin
//go:generate yaegi extract "github.com/kuchensheng/bintools/tracer/trace"
//go:generate yaegi extract "github.com/rs/zerolog/log"

package extractlib

import "reflect"

// Symbols variable stores the map of stdlib symbols per package.
var Symbols = map[string]map[string]reflect.Value{}

func init() {
	Symbols["github.com/traefik/yaegi/stdlib/stdlib"] = map[string]reflect.Value{
		"Symbols": reflect.ValueOf(Symbols),
	}
}
