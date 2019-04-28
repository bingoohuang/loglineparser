# logfmtparser
log parser to golang struct

## 日志格式定义

```text
2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [200, 0.023999929428101, 1539866805.135, 108],  [true, -, -], [{}] request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"

```

在形如以上日志行的日志中，以[xx]包含起来的，是需要提取的部分(parts)，索引号从0开始。

然后在xx中，可能有多个子字段，比如[200, 0.023999929428101, 1539866805.135, 108]，这个以逗号分隔的，是自字段。

可以定义以下go语言的结构体，来映射这些提取部分，或者提取子字段(subs):

```go
type LogLine struct {
	LogLevel       string `llp:"0" json:"logLevel"` // notice
	GatewayStatus string `llp:"2" json:"gatewayFlag"` // GatewayMonV2

	RespStatus            string    `llp:"4.0" json:"respStatus"`
	RespResponseTime      float32   `llp:"4.1" json:"respResponseTime"`
	RespInnerStartReqTime time.Time `llp:"4.2" json:"respInnerStartReqTime"`
	RespBodySize          int       `llp:"4.3" json:"respBodySize"`


	AuthIsLocalIP          bool   `llp:"5.0" json:"authIsLocalIP"`
	AuthKeySecretCheckRst  string `llp:"5.1" json:"authKeySecretCheckRst"`
	
	ApiSessionVarMap map[string]string `llp:"6" json:"apiSessionVarMap"`
}


// ParseLogLine 解析一行日志
func ParseLogLine(line string) (LogLine, error) {
	v := LogLine{}
	err := logfmtparser.ParseLogLine(line, &v)

	return v, err
}


```

其中，结构体LogLine各个字段tag中的`llp`（loglineparser的缩写）部分，使用x或者x.y的表达方式：

1. x 表示取第x个（从0开始）提取值，并且根据需要，进行合适的类型转换。
1. x.y 表示取第x个（从0开始）提取值的第y个（从0开始）子值，进行合适的类型转换。


如果，需要实现自定义解码，可以参考以下示例：

```go
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

```


