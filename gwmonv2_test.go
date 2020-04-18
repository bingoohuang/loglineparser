// nolint lll
package loglineparser_test

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
)

var gatewayMonV2Parser = loglineparser.NewLogLineParser((*loglineparser.GatewayMonV2)(nil))

// or var gatewayMonV2Parser = loglineparser.NewLogLineParser("loglineparser.GatewayMonV2")

// ParseGatewayMonV2 解析GatewayMonV2日志
func ParseGatewayMonV2(line string) (loglineparser.GatewayMonV2, error) {
	v, err := gatewayMonV2Parser.Parse(line)

	return v.(loglineparser.GatewayMonV2), err
}

func BenchmarkParseGatewayMonV2(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023999929428101, 0.023999929428101, 1539866805.135, 108], [1539866805.135, 192.168.106.8, -, 208],` +
			` [-, 127.0.0.1-1539866805.135-19002-2879, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081`
		ParseGatewayMonV2(line)
	}
}

func BenchmarkFastParseGatewayMonV2(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023999929428101, 0.023999929428101, 1539866805.135, 108], [1539866805.135, 192.168.106.8, -, 208],` +
			` [-, 127.0.0.1-1539866805.135-19002-2879, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081`
		loglineparser.FastCreateGatewayMonV2(line)
	}
}

func TestParseGatewayMonV2(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023999929428101, 0.023999929428101, 1539866805.135, 108], [1539866805.135, 192.168.106.8, -, 208],` +
		` [-, 127.0.0.1-1539866805.135-19002-2879, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081`
	v, err := ParseGatewayMonV2(line)
	v2 := loglineparser.FastCreateGatewayMonV2(line)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(v2, v)
	a.Equal(loglineparser.GatewayMonV2{
		LogType:       "GatewayMonV2",
		GatewayStatus: "200",
		APIVersionID:  "10031",
		RespStatus:    "200",

		RespResponseTime:      0.023999929428101,
		RespInnerResponseTime: 0.023999929428101,
		RespInnerStartReqTime: loglineparser.ParseTime("1539866805.135"),
		RespBodySize:          108,

		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: "192.168.106.8",
		UserUID:      "",
		RequestSize:  208,

		RequestID: "",
		TraceID:   "127.0.0.1-1539866805.135-19002-2879",
		ServiceID: "",

		AuthIsLocalIP:          true,
		AuthKeySecretCheckRst:  "",
		AuthSessionCheckRst:    "",
		AuthSid:                "",
		AuthUcenterPlatform:    "",
		AuthCheckLoginTokenRst: "",
		AuthCookie:             "",
		AuthInvalidMsg:         "",

		APISessionVarMap: map[string]string{},

		UserRealIP:        "",
		UserUa:            "Apache-HttpClient/4.5.5 (Java/1.8.0_181)",
		LastHop:           "192.168.106.8",
		UserXForwardedFor: "",
	}, v)
}

type MyIP net.IP

func (i *MyIP) Unmarshal(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return errors.New("bad format ip " + v)
	}
	*i = MyIP(ip)

	return nil
}

var _ loglineparser.Unmarshaler = (*MyIP)(nil)

// 实现参考自: https://github.com/projectcalico/libcalico-go/blob/master/lib/net/ip.go
func (i MyIP) MarshalJSON() ([]byte, error) {
	s, err := net.IP(i).MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(s))
}

type LogLine struct {
	LogType string `llp:"2" json:"logType"` // GatewayMonV2

	UserTime     time.Time `llp:"3.0" json:"reqTime"`
	UserClientIP MyIP      `llp:"3.1" json:"userClientIP"`

	Xy string `llp:"4" json:"xy"`
}

func TestCustomDecode(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2], [1539866805.135, 192.168.106.8, -, 208] [x,y] xxxxx`

	var LogLineParser = loglineparser.NewLogLineParser(LogLine{})

	v, err := LogLineParser.Parse(line)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(LogLine{
		LogType:      "GatewayMonV2",
		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: MyIP(net.ParseIP("192.168.106.8")),
		Xy:           "x,y",
	}, v)
}

type LogLineUser struct {
	UserTime     time.Time `llp:"3.0" json:"reqTime"`
	UserClientIP MyIP      `llp:"3.1" json:"userClientIP"`
}

type LogLine2 struct {
	LogType string `llp:"2" json:"logType"` // GatewayMonV2

	LogLineUser

	Xy string `llp:"4" json:"xy"`
}

var LogLine2Parser = loglineparser.NewLogLineParser("loglineparser_test.LogLine2")

func TestCustomDecode2(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2], [1539866805.135, 192.168.106.8, -, 208] [x,y] xxxxx`

	v, err := LogLine2Parser.Parse(line)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(LogLine2{
		LogType: "GatewayMonV2",
		LogLineUser: LogLineUser{UserTime: loglineparser.ParseTime("1539866805.135"),
			UserClientIP: MyIP(net.ParseIP("192.168.106.8"))},
		Xy: "x,y",
	}, v)
}
