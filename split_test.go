package loglineparser_test

import (
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkMakeBracketPartSplitter(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := createParser()
		for {
			parser.ParseParts()
		}
	}
}

func createParser() loglineparser.PartSplitter {
	parser := loglineparser.MakeBracketPartSplitter()
	parser.LoadLine(
		`2019/04/27 03:12:01 [notice] 17618#0: *579576278 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [2001], [-], ` +
			`[404, -, -, -, -], [1556305921.3, 100.120.36.178, -, 19], [-, 127.0.0.1-1556305921.3-17618-470, -], [false, -, -, -, -, -, -, -], [-], ` +
			`[-, -, 100.120.36.178, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-], client: 100.120.36.178, server: localhost, request: "HEAD / HTTP/1.0"`)
	return parser
}

func TestMakeBracketPartSplitter(t *testing.T) {
	a := assert.New(t)

	parser := createParser()

	parts := []string{"notice", "lua", "GatewayMonV2", "2001", "-", "404, -, -, -, -", "1556305921.3, 100.120.36.178, -, 19", "-, 127.0.0.1-1556305921.3-17618-470, -",
		"false, -, -, -, -, -, -, -", "-", "-, -, 100.120.36.178, -", "-", "-", "-, -, -, -, -, -, -, -, -, -, -", "-End-"}

	a.Equal(parts, parser.ParseParts())
}

func TestMakeBracketPartSplitterEx(t *testing.T) {
	a := assert.New(t)

	parser := loglineparser.MakeBracketPartSplitter()
	parser.LoadLine(
		`2019/04/27 03:12:01 [not[ic]e] 17618#0: *579576278 [lua] gateway.lua:163: log_base(): [GatewayMonV2][2001], [-], `)

	parts := []string{"not[ic]e", "lua", "GatewayMonV2", "2001", "-"}

	a.Equal(parts, parser.ParseParts())
}
