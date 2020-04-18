package loglineparser_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
)

func TestIsBlank(t *testing.T) {
	a := assert.New(t)
	a.True(loglineparser.IsBlank(""))
	a.True(loglineparser.IsBlank(" "))
	a.True(loglineparser.IsBlank("　"))
	a.True(loglineparser.IsBlank("\t\r\n"))

	type Papa struct {
		Name string
	}

	pv := reflect.ValueOf((*Papa)(nil))
	fmt.Println(pv.Type())
	elem := pv.Type().Elem()
	ev := reflect.Zero(elem).Interface().(Papa)
	ev.Name = "bingoo"
}
