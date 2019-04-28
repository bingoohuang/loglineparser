package loglineparser_test

import (
	"errors"
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestParseGatewayMonV2(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023999929428101, 0.023999929428101, 1539866805.135, 108], [1539866805.135, 192.168.106.8, -, 208],` +
		` [-, 127.0.0.1-1539866805.135-19002-2879, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081`
	v, err := loglineparser.ParseGatewayMonV2(line)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(loglineparser.GatewayMonV2{
		LogType:       "GatewayMonV2",
		GatewayStatus: "200",
		ApiVersionId:  "10031",
		RespStatus:    "200",

		RespResponseTime:      0.023999929428101,
		RespInnerResponseTime: 0.023999929428101,
		RespInnerStartReqTime: loglineparser.ParseTime("1539866805.135"),
		RespBodySize:          108,

		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: "192.168.106.8",
		UserUid:      "",
		RequestSize:  208,

		RequestId: "",
		TraceId:   "127.0.0.1-1539866805.135-19002-2879",
		ServiceId: "",

		AuthIsLocalIP:          true,
		AuthKeySecretCheckRst:  "",
		AuthSessionCheckRst:    "",
		AuthSid:                "",
		AuthUcenterPlatform:    "",
		AuthCheckLoginTokenRst: "",
		AuthCookie:             "",
		AuthInvalidMsg:         "",

		ApiSessionVarMap: map[string]string{},

		UserRealIP:        "",
		UserUa:            "Apache-HttpClient/4.5.5 (Java/1.8.0_181)",
		LastHop:           "192.168.106.8",
		UserXForwardedFor: "",
	}, v)
}

type IP struct {
	ip net.IP
}

func MakeIP(v string) IP {
	return IP{ip: net.ParseIP(v)}
}

func (i *IP) Unmarshal(v string) error {
	i.ip = net.ParseIP(v)
	if i.ip == nil {
		return errors.New("bad format ip " + v)
	}

	return nil
}

func (i *IP) MarshalJSON() ([]byte, error) {
	return i.MarshalJSON()
}

type LogLine struct {
	LogType string `llp:"2" json:"logType"` // GatewayMonV2

	UserTime     time.Time `llp:"3.0" json:"reqTime"`
	UserClientIP IP        `llp:"3.1" json:"userClientIP"`
}

func TestCustomDecode(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2], [1539866805.135, 192.168.106.8, -, 208] xxxxx`
	v := LogLine{}
	err := loglineparser.ParseLogLine(line, &v)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(LogLine{
		LogType:      "GatewayMonV2",
		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: MakeIP("192.168.106.8"),
	}, v)
}
