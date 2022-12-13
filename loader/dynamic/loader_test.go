package dynamic

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"testing"
)

func TestLoad(t *testing.T) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	v, err := i.Eval(`import "fmt"`)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Kind().String(), v.Interface())

	v, err = i.Eval(`fmt.Println("Hello Yaegi")`)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Kind().String(), v.Interface())
}
