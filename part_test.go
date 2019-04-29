package loglineparser_test

import (
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkMakeBracketPartSplitter(b *testing.B) {
	b.ResetTimer()
	parser := loglineparser.NewBracketPartSplitter()
	for i := 0; i < b.N; i++ {
		parseLine(parser)
	}
}

func parseLine(parser loglineparser.PartSplitter) []string {
	return parser.Parse(
		`2019/04/27 03:12:01 [notice] 17618#0: *579576278 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [2001], [-], ` +
			`[404, -, -, -, -], [1556305921.3, 100.120.36.178, -, 19], [-, 127.0.0.1-1556305921.3-17618-470, -], [false, -, -, -, -, -, -, -], [-], ` +
			`[-, -, 100.120.36.178, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-], client: 100.120.36.178, server: localhost, request: "HEAD / HTTP/1.0"`)
}

func TestMakeBracketPartSplitter(t *testing.T) {
	a := assert.New(t)
	parser := loglineparser.NewBracketPartSplitter()

	parts := []string{"notice", "lua", "GatewayMonV2", "2001", "-", "404, -, -, -, -", "1556305921.3, 100.120.36.178, -, 19", "-, 127.0.0.1-1556305921.3-17618-470, -",
		"false, -, -, -, -, -, -, -", "-", "-, -, 100.120.36.178, -", "-", "-", "-, -, -, -, -, -, -, -, -, -, -", "-End-"}

	a.Equal(parts, parseLine(parser))
}

func TestMakeBracketPartSplitterEx(t *testing.T) {
	a := assert.New(t)

	parser := loglineparser.NewBracketPartSplitter()
	realParts := parser.Parse(
		`2019/04/27 03:12:01 [not[ic]e] 17618#0: *579576278 [lua] gateway.lua:163: log_base(): [GatewayMonV2][2001], [-], `)

	parts := []string{"not[ic]e", "lua", "GatewayMonV2", "2001", "-"}

	a.Equal(parts, realParts)
}
