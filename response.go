package mute

import (
	"encoding/json"
	"net/http"
)

// muteHttpResponse 响应
// response http.Response
// client mute http 客户端
// body 响应体
type muteHttpResponse struct {
	response *http.Response
	client   muteHttpClient
	body     []byte
}

// Code http code
func (r *muteHttpResponse) Code() int {
	if r.response == nil {
		return 0
	}
	return r.response.StatusCode
}

// GetBody body
func (r *muteHttpResponse) GetBody() []byte {
	return r.body
}

// Cookies response Cookies
func (r *muteHttpResponse) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

// Request http request
func (r *muteHttpResponse) Request() *http.Request {
	return r.client.request
}

// Curl curl cmd string
func (r *muteHttpResponse) Curl() string {
	if r.client.request == nil {
		return ""
	}
	return buildCurl(r.client.url, r.client.method, string(r.client.body), r.client.request.Header)
}

// UseTime 请求耗时
func (r *muteHttpResponse) UseTime() int64 {
	return r.client.useTime
}

// Unmarshal JSON反序列化
func (r *muteHttpResponse) Unmarshal(resp interface{}) error {
	return json.Unmarshal(r.body, resp)
}
