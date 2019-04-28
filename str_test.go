package loglineparser_test

import (
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBlank(t *testing.T) {
	a := assert.New(t)
	a.True(loglineparser.IsBlank(""))
	a.True(loglineparser.IsBlank(" "))
	a.True(loglineparser.IsBlank("　"))
	a.True(loglineparser.IsBlank("\t\r\n"))
}
