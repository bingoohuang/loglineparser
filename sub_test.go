package loglineparser_test

import (
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeLogPartSplitter(t *testing.T) {
	a := assert.New(t)

	parser := loglineparser.MakeSubSplitter(",", "-")
	a.Equal([]string{"", "", "100.120.36.178", ""}, parser.Parse("-, -, 100.120.36.178, -"))
}
