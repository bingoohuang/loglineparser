// nolint lll
package loglineparser_test

import (
	"testing"

	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
)

func TestParseNegativeOffset(t *testing.T) {
	type app struct {
		LogLevel   string `llp:"0"`     // notice，提取第0+1个[]中整个部分
		AppID      string `llp:"10.-3"` // 应用 ID
		CustomerID string `llp:"10.-2"` // 客户 ID
		DeviceID   string `llp:"10.-1"` // 设备 ID
	}

	parser, err := loglineparser.New(app{})
	assert.Nil(t, err)

	line := `2022/12/19 11:04:30 [notice] 15639#0: *1339077289 [lua] gw_log.lua:121: log_base(): [GatewayMonV2] [200], ` +
		`[30427], [200, 0.57399988174438, 0.54399991035461, 1671419070.372, 52145], [1671419070.342, 61.155.4.121, -, 446], ` +
		`[-, 127.0.0.1-1671419070.342-15639-681, -], [false, true, -, -, -, -, -, -], [{}], ` +
		`[-, MSSP-User-Agent, 100.120.36.16, 168.188.9.42, 61.155.4.121, 0.00399, 52216, APP_E898287, CST_E932B93, DEV_0A83AB6], ` +
		`[-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 61.155.4.121, , ` +
		`server: localhost, request: "POST /idsctid/v2a/api/identities/mode0x42 HTTP/1.0", host: "168.188.3.11"`
	ap, err := parser.Parse(line)
	assert.Nil(t, err)
	assert.Equal(t, app{
		LogLevel:   "notice",
		AppID:      "APP_E898287",
		CustomerID: "CST_E932B93",
		DeviceID:   "DEV_0A83AB6",
	}, ap.(app))
}

func BenchmarkMakeBracketPartSplitter(b *testing.B) {
	b.ResetTimer()

	parser := loglineparser.NewBracketPartSplitter("-")

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
	parser := loglineparser.NewBracketPartSplitter("-")

	parts := []string{
		"notice", "lua", "GatewayMonV2", "2001", "", "404, -, -, -, -", "1556305921.3, 100.120.36.178, -, 19", "-, 127.0.0.1-1556305921.3-17618-470, -",
		"false, -, -, -, -, -, -, -", "", "-, -, 100.120.36.178, -", "", "", "-, -, -, -, -, -, -, -, -, -, -", "-End-",
	}

	a.Equal(parts, parseLine(parser))
}

func TestMakeBracketPartSplitterEx(t *testing.T) {
	a := assert.New(t)

	parser := loglineparser.NewBracketPartSplitter("-")
	realParts := parser.Parse(
		`2019/04/27 03:12:01 [not[ic]e] 17618#0: *579576278 [lua] gateway.lua:163: log_base(): [GatewayMonV2][2001], [-], `)

	parts := []string{"not[ic]e", "lua", "GatewayMonV2", "2001", ""}

	a.Equal(parts, realParts)
}
