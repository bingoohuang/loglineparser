package loglineparser_test

import (
	"testing"

	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
)

func TestMakeLogPartSplitter(t *testing.T) {
	a := assert.New(t)

	parser := loglineparser.NewSubSplitter(",", "-")
	a.Equal([]string{"", "", "100.120.36.178", ""}, parser.Parse("-, -, 100.120.36.178, -"))
}
