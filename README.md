# loglineparser

[![Travis CI](https://img.shields.io/travis/bingoohuang/loglineparser/master.svg?style=flat-square)](https://travis-ci.com/bingoohuang/loglineparser)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/bingoohuang/loglineparser/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/bingoohuang/loglineparser)
[![Coverage Status](http://codecov.io/github/bingoohuang/loglineparser/coverage.svg?branch=master)](http://codecov.io/github/bingoohuang/loglineparser?branch=master)
[![goreport](https://www.goreportcard.com/badge/github.com/bingoohuang/loglineparser)](https://www.goreportcard.com/report/github.com/bingoohuang/loglineparser)

log parser to parse log line to relative golang struct.

## 日志格式定义

```text
2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [200, 0.023999929428101, 1539866805.135, 108],  [true, -, -], [{}] request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
```

1. 在形如以上日志行的日志中，以[xx]包含起来的，是需要提取的部分(parts)，索引号从0开始。 然后在xx中，可能有多个子字段(subs)，比如[200, 0.023999929428101, 1539866805.135, 108]，这个以逗号分隔的，是子字段(subs)。
2. 或者使用正则表达式的捕获组提取字段值，见下面的示例

可以定义以下go语言的结构体，来映射这些提取部分(parts)，或者提取子字段(subs):

```go
package yourawesomepackage

import (
	"github.com/bingoohuang/loglineparser"
	"time"
)

type LogLine struct {
	LogLevel      string `llp:"0"`    // notice
	GatewayStatus string `llp:"2"` // GatewayMonV2

	RespStatus            string    `llp:"4.0"`
	RespResponseTime      float32   `llp:"4.1"`
	RespInnerStartReqTime time.Time `llp:"4.2"`
	RespBodySize          int       `llp:"4.3"`


	AuthIsLocalIP          bool   `llp:"5.0"`
	AuthKeySecretCheckRst  string `llp:"5.1"`
	
	ApiSessionVarMap map[string]string `llp:"6"`


	// 例如：捕获 request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1" 中的 POST 部分
	Method string `llp:"reg" reg:"request: \"([A-Z]+)"` // 使用正则表达式捕获，默认捕获组 1
	// 例如：捕获 host: "192.168.108.11:8081" 中的 192.168.108.11:8081 部分
	Host string `llp:"reg" reg:"host: \"(.*?)\"" group:"1"` // 使用正则表达式捕获，捕获组序号 1
}


var logLineParser, _ = loglineparser.New(&LogLine{})

// ParseLogLine 解析一行日志
func ParseLogLine(line string) (*LogLine, error) {
	v, err := logLineParser.Parse(line)
    if err != nil {
        return nil, err
    }

	return v.(*LogLine), nil
}
```

其中，结构体LogLine各个字段tag中的`llp`（loglineparser的缩写）部分，使用以下表达方式：

1. x 表示取第x个（从0开始）提取值，并且根据需要，进行合适的类型转换。
1. x.y 表示取第x个（从0开始）提取值的第y个（从0开始）子值，进行合适的类型转换。

如果需要实现自定义解码，可以参考以下示例：

```go
package yourawesomepackage

import (
	"encoding/json"
	"errors"
	"github.com/bingoohuang/loglineparser"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type MyIP net.IP

func (i *MyIP) Unmarshal(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return errors.New("bad format ip " + v)
	}
	*i = MyIP(ip)

	return nil
}

type LogLine struct {
	LogType      string     `llp:"2"`
	UserTime     time.Time  `llp:"3.0"`
	UserClientIP MyIP       `llp:"3.1"`
}


var logLineParser, _ = loglineparser.New(&LogLine{})

func TestCustomDecode(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2], [1539866805.135, 192.168.106.8, -, 208] [x,y] xxxxx`

	v, err := logLineParser.Parse(line)

	assert.Nil(t, err)
	assert.Equal(t, LogLine{
		LogType:      "GatewayMonV2",
		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: MyIP(net.ParseIP("192.168.106.8")),
	}, *v)
}
```


## 运行测试

1. 运行测试用例 `go fmt ./...; go test ./... -v -count=1`
1. 运行基准用例 `go test -bench=.`

```bash
$ go test ./...
ok  	github.com/bingoohuang/loglineparser	0.013s

$ go test -bench=.
*loglineparser_test.Papa
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/loglineparser
BenchmarkParseGatewayMonV2-12          	   20000	     97816 ns/op
BenchmarkFastParseGatewayMonV2-12      	   20000	     81113 ns/op
BenchmarkMakeBracketPartSplitter-12    	   30000	     52485 ns/op
PASS
ok  	github.com/bingoohuang/loglineparser	4.730s
```

> [97816 Nanoseconds = 0.097816 Milliseconds](https://convertlive.com/u/convert/nanoseconds/to/milliseconds#97816)
