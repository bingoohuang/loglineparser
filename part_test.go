package loglineparser_test

import (
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeLogPartSplitter(t *testing.T) {
	a := assert.New(t)

	sp := loglineparser.MakeSubSplitter()
	sp.Load("-, -, 100.120.36.178, -", ",")
	a.Equal(4, sp.Len())
	a.Equal("", sp.Sub(0))
	a.Equal("", sp.Sub(1))
	a.Equal("100.120.36.178", sp.Sub(2))
	a.Equal("", sp.Sub(3))
	a.Equal("", sp.Sub(4))

	a.Equal([]string{"", "", "100.120.36.178", ""}, sp.Subs())
}
