package loglineparser

import (
	"time"
)

// ParseGatewayMonV2 解析GatewayMonV2日志
func ParseGatewayMonV2(line string) (GatewayMonV2, error) {
	v := GatewayMonV2{}
	err := ParseLogLine(line, &v)

	return v, err
}

/*
 -- http://192.168.131.32:9000/develop/FOOTSTONE/GateWay/Code/api-gateway-ng/blob/master/lua/web/gateway.lua
 ngx.log(ngx.NOTICE, string.format("[%s] [%s], [%s], [%s, %s, %s, %s, %s], [%s, %s, %s, %s], [%s, %s, %s], [%s, %s, %s, %s, %s, %s, %s, %s], [%s], [%s, %s, %s, %s], [%s], [%s], [%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s], [-End-]",
            log_format_name,
            rec_tab.flag,
            rec_tab.api_apiid,
            rec_tab.resp_status, rec_tab.resp_response_time, rec_tab.resp_inner_response_time, rec_tab.resp_inner_start_req_time, rec_tab.resp_body_size,
            rec_tab.user_time, rec_tab.user_client_ip, rec_tab.user_uid, rec_tab.request_size,
            rec_tab.request_id, rec_tab.trace_id, rec_tab.service_id,
            rec_tab.auth_is_local_ip, rec_tab.auth_key_secret_check_rst, rec_tab.auth_session_check_rst, rec_tab.auth_sid, rec_tab.auth_ucenter_platform, rec_tab.auth_check_login_token_rst, rec_tab.auth_cookie, rec_tab.auth_invalid_msg,
            rec_tab.api_session_var_map,
            rec_tab.user_realip, rec_tab.user_ua, rec_tab.last_hop, rec_tab.user_x_forwarded_for,
            rec_tab.request_body,
            rec_tab.resp_body,
            rec_tab.anchor_init_var, rec_tab.anchor_search_api,
            rec_tab.anchor_check_develop, rec_tab.anchor_check_app,rec_tab.anchor_check_session,
            rec_tab.anchor_check_safe,
            rec_tab.anchor_start_request,
            rec_tab.anchor_produce_cookie,rec_tab.anchor_encrypt_cookie,
            rec_tab.anchor_produce_response,
            rec_tab.anchor_output_log))

2018/10/18 20:46:45 [notice] 19002#0: *53103423 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023999929428101, 0.023999929428101, 1539866805.135, 108], [1539866805.135, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.135-19002-2879, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19001#0: *53103910 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.013000011444092, 0.013000011444092, 1539866805.146, 108], [1539866805.146, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.146-19001-3366, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19002#0: *53103729 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.027000188827515, 0.027000188827515, 1539866805.133, 108], [1539866805.133, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.133-19002-3185, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19002#0: *53100709 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.30299997329712, 0.30299997329712, 1539866804.858, 108], [1539866804.858, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866804.858-19002-165, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 18999#0: *53103890 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.023000001907349, 0.023000001907349, 1539866805.138, 108], [1539866805.138, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.138-18999-3346, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19000#0: *53100582 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.013000011444092, 0.013000011444092, 1539866805.149, 108], [1539866805.149, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.149-19000-38, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19002#0: *53100550 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.021000146865845, 0.021000146865845, 1539866805.142, 108], [1539866805.142, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.142-19002-6, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19002#0: *53103874 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.032000064849854, 0.032000064849854, 1539866805.132, 108], [1539866805.132, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.132-19002-3330, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 18999#0: *53100390 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.30400013923645, 0.30400013923645, 1539866804.86, 108], [1539866804.86, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866804.86-18999-3942, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.11:8081"
2018/10/18 20:46:45 [notice] 19000#0: *53100587 [lua] gateway.lua:163: log_base(): [GatewayMonV2] [200], [10031], [200, 0.024999856948853, 0.024999856948853, 1539866805.14, 108], [1539866805.14, 192.168.106.8, -, 208], [-, 127.0.0.1-1539866805.14-19000-43, -], [true, -, -, -, -, -, -, -], [{}], [-, Apache-HttpClient/4.5.5 (Java/1.8.0_181), 192.168.106.8, -], [-], [-], [-, -, -, -, -, -, -, -, -, -, -], [-End-] while sending to client, client: 192.168.106.8, server: localhost, request: "POST /dsvs/v1/pkcs1/verifyDigestSign HTTP/1.1", host: "192.168.108.1
*/
type GatewayMonV2 struct {
	LogType       string `llp:"2" json:"logType"` // GatewayMonV2
	GatewayStatus string `llp:"3" json:"gatewayFlag"`
	ApiVersionId  string `llp:"4" json:"apiVersionId"`

	RespStatus            string    `llp:"5.0" json:"respStatus"`
	RespResponseTime      float32   `llp:"5.1" json:"respResponseTime"`
	RespInnerResponseTime float32   `llp:"5.2" json:"respInnerResponseTime"`
	RespInnerStartReqTime time.Time `llp:"5.3" json:"respInnerStartReqTime"`
	RespBodySize          int       `llp:"5.4" json:"respBodySize"`

	UserTime     time.Time `llp:"6.0" json:"reqTime"`
	UserClientIP string    `llp:"6.1" json:"userClientIP"`
	UserUid      string    `llp:"6.2" json:"userUid"`
	RequestSize  int       `llp:"6.3" json:"requestSize"`

	RequestId string `llp:"7.0" json:"requestId"`
	TraceId   string `llp:"7.1" json:"TraceId"`
	ServiceId string `llp:"7.2" json:"serviceId"`

	AuthIsLocalIP          bool   `llp:"8.0" json:"authIsLocalIP"`
	AuthKeySecretCheckRst  string `llp:"8.1" json:"authKeySecretCheckRst"`
	AuthSessionCheckRst    string `llp:"8.2" json:"authSessionCheckRst"`
	AuthSid                string `llp:"8.3" json:"authSid"`
	AuthUcenterPlatform    string `llp:"8.4" json:"authUcenterPlatform"`
	AuthCheckLoginTokenRst string `llp:"8.5" json:"authCheckLoginTokenRst"`
	AuthCookie             string `llp:"8.6" json:"authCookie"`
	AuthInvalidMsg         string `llp:"8.7" json:"authInvalidMsg"`

	ApiSessionVarMap map[string]string `llp:"9" json:"apiSessionVarMap"`

	UserRealIP        string `llp:"10.0" json:"userRealIP"`
	UserUa            string `llp:"10.1" json:"userUa"`
	LastHop           string `llp:"10.2" json:"lastHop"`
	UserXForwardedFor string `llp:"10.3" json:"userXForwardedFor"`

	RequestBody string `llp:"11" json:"requestBody"`
	RespBody    string `llp:"12" json:"respBody"`

	AnchorInitVar         string `llp:"13.0" json:"anchorInitVar"`
	AnchorSearchApi       string `llp:"13.1" json:"anchorSearchApi"`
	AnchorCheckDevelop    string `llp:"13.2" json:"anchorCheckDevelop"`
	AnchorCheckApp        string `llp:"13.3" json:"anchorCheckApp"`
	AnchorCheckSession    string `llp:"13.4" json:"anchorCheckSession"`
	AnchorCheckSafe       string `llp:"13.5" json:"anchorCheckSafe"`
	AnchorStartRequest    string `llp:"13.6" json:"anchorStartRequest"`
	AnchorProduceCookie   string `llp:"13.7" json:"anchorProduceCookie"`
	AnchorEncryptCookie   string `llp:"13.8" json:"anchorEncryptCookie"`
	AnchorProduceResponse string `llp:"13.9" json:"anchorProduceResponse"`
	AnchorOutputLog       string `llp:"13.10" json:"anchorOutputLog"`
}
