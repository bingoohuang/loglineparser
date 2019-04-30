# loglineparser
log parser to golang struct

## 日志格式定义

```text
2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [200, 0.023999929428101, 1539866805.135, 108],  [true, -, -], [{}] request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"

```

在形如以上日志行的日志中，以[xx]包含起来的，是需要提取的部分(parts)，索引号从0开始。

然后在xx中，可能有多个子字段(subs)，比如[200, 0.023999929428101, 1539866805.135, 108]，这个以逗号分隔的，是子字段(subs)。

可以定义以下go语言的结构体，来映射这些提取部分，或者提取子字段(subs):

```go
package yourawesomepackage

import (
	"github.com/bingoohuang/loglineparser"
	"time"
)

type LogLine struct {
	LogLevel      string `llp:"0" json:"logLevel"`    // notice
	GatewayStatus string `llp:"2" json:"gatewayFlag"` // GatewayMonV2

	RespStatus            string    `llp:"4.0" json:"respStatus"`
	RespResponseTime      float32   `llp:"4.1" json:"respResponseTime"`
	RespInnerStartReqTime time.Time `llp:"4.2" json:"respInnerStartReqTime"`
	RespBodySize          int       `llp:"4.3" json:"respBodySize"`


	AuthIsLocalIP          bool   `llp:"5.0" json:"authIsLocalIP"`
	AuthKeySecretCheckRst  string `llp:"5.1" json:"authKeySecretCheckRst"`
	
	ApiSessionVarMap map[string]string `llp:"6" json:"apiSessionVarMap"`
}


var LogLineParser = loglineparser.NewLogLineParser("yourawesomepackage.LogLine")

// ParseLogLine 解析一行日志
func ParseLogLine(line string) (LogLine, error) {
	v, err := LogLineParser.Parse(line)
	return v.(LogLine), err
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

type MyIP struct {
	net.IP
}

func (i *MyIP) Unmarshal(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return errors.New("bad format ip " + v)
	}
	i.IP = ip

	return nil
}

func (i MyIP) MarshalJSON() ([]byte, error) {
	s, err := i.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(s))
}

type LogLine struct {
	LogType string `llp:"2" json:"logType"` // GatewayMonV2

	UserTime     time.Time `llp:"3.0" json:"reqTime"`
	UserClientIP MyIP       `llp:"3.1" json:"userClientIP"`
}


var LogLineParser = loglineparser.NewLogLineParser("yourawesomepackage.LogLine")

func TestCustomDecode(t *testing.T) {
	line := `2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2], [1539866805.135, 192.168.106.8, -, 208] [x,y] xxxxx`

	v, err := LogLineParser.Parse(line)

	a := assert.New(t)
	a.Nil(err)
	a.Equal(LogLine{
		LogType:      "GatewayMonV2",
		UserTime:     loglineparser.ParseTime("1539866805.135"),
		UserClientIP: MyIP{net.ParseIP("192.168.106.8")},
	}, v)
}


```


## 运行测试

1. 运行测试用例 `go test ./...`
1. 运行基准用例 `go test -bench=.`

```bash
$ go test ./...
ok  	github.com/bingoohuang/loglineparser	(cached)

$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/loglineparser
BenchmarkParseGatewayMonV2-12          	   10000	    102426 ns/op
BenchmarkMakeBracketPartSplitter-12    	   30000	     51495 ns/op
PASS
ok  	github.com/bingoohuang/loglineparser	3.138s
```

> [102426 Nanoseconds = 0.102426 Milliseconds](https://convertlive.com/u/convert/nanoseconds/to/milliseconds#102426)
