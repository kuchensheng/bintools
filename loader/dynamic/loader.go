package dynamic

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var scriptEngine = func() *interp.Interpreter {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	return i
}()

func LoadGoScript(script string) (res any, err error) {

	eval, err := scriptEngine.Eval(script)
	if err != nil {
		return nil, err
	}
	return eval.Interface(), nil
}
