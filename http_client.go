package mute

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// codeMustErr code must error
	codeMustErr = errors.New("HTTP code mismatch")

	// defaultClient default http client
	defaultClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				c, err := d.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
)

// muteHttpClient 请求client
// url 请求url
// mustCode 限定的返回http status code
// method 请求方法
// body 请求body
// request http.Request
// client http.Client
// useTime 请求耗时
type muteHttpClient struct {
	url      string
	mustCode int
	method   string
	body     []byte
	request  *http.Request
	client   *http.Client
	useTime  int64
}

// NewWithClient create muteHttpClient with you http client
func NewWithClient(url string, client *http.Client) *muteHttpClient {
	return &muteHttpClient{url: url, request: &http.Request{Header: make(http.Header, 16)}, client: client}
}

// NewWithTransport create muteHttpClient with you http transport
func NewWithTransport(url string, transport *http.Transport) *muteHttpClient {
	client := defaultClient
	client.Transport = transport
	return NewWithClient(url, client)
}

// New create muteHttpClient with default
func New(url string) *muteHttpClient {
	return NewWithClient(url, defaultClient)
}

// SetBodyJSON set json body
func (c *muteHttpClient) SetBodyJSON(obj interface{}) *muteHttpClient {
	c.body, _ = json.Marshal(obj)
	c.request.Header["Content-Type"] = []string{"application/json"}
	c.request.Body = io.NopCloser(bytes.NewReader(c.body))
	return c
}

// AddCookie add Cookie
func (c *muteHttpClient) AddCookie(cookies ...*http.Cookie) *muteHttpClient {
	for _, cookie := range cookies {
		c.request.AddCookie(cookie)
	}
	return c
}

// AddHeader append http header
func (c *muteHttpClient) AddHeader(key, value string) *muteHttpClient {
	if _, ok := c.request.Header[key]; ok {
		c.request.Header[key] = append(c.request.Header[key], value)
	} else {
		c.request.Header[key] = []string{value}
	}
	return c
}

// SetHeader set http header
func (c *muteHttpClient) SetHeader(key, value string) *muteHttpClient {
	c.request.Header[key] = []string{value}
	return c
}

// SetQuery set query param
func (c *muteHttpClient) SetQuery(key, value string) *muteHttpClient {
	c.request.URL, _ = url.Parse(c.url)
	values, _ := url.ParseQuery(c.request.URL.RawQuery)
	if values == nil {
		values = make(url.Values)
	}
	values.Set(key, value)
	jointStr := "?"
	if len(values) > 1 {
		jointStr = ""
	}
	c.url = strings.TrimSuffix(c.url, c.request.URL.RawQuery) + jointStr + values.Encode()
	return c
}

// Header append http header
func (c *muteHttpClient) Header(header http.Header) *muteHttpClient {
	if len(header) == 0 {
		c.request.Header = make(http.Header)
	} else {
		c.request.Header = header
	}
	return c
}

// SetPostForm set post form
func (c *muteHttpClient) SetPostForm(value url.Values) *muteHttpClient {
	c.request.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	c.request.PostForm = value
	c.body = []byte(value.Encode())
	return c
}

// SetCookieJar set CookieJar
func (c *muteHttpClient) SetCookieJar(cookieJar http.CookieJar) *muteHttpClient {
	if c.client == defaultClient {
		c.dumpDefaultClient()
	}
	c.client.Jar = cookieJar
	return c
}

// dumpDefaultClient dump with default client
func (c *muteHttpClient) dumpDefaultClient() {
	client := *defaultClient
	c.client = &client
}

// MustCode set mustCode
func (c *muteHttpClient) MustCode(code int) *muteHttpClient {
	c.mustCode = code
	return c
}

// Post do
func (c *muteHttpClient) Post(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodPost, ctx)
}

// Get do
func (c *muteHttpClient) Get(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodGet, ctx)
}

// Put do
func (c *muteHttpClient) Put(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodPut, ctx)
}

// Delete do
func (c *muteHttpClient) Delete(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodDelete, ctx)
}

// Options do
func (c *muteHttpClient) Options(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodOptions, ctx)
}

// Patch do
func (c *muteHttpClient) Patch(ctx context.Context) (muteHttpResponse, error) {
	return c.do(http.MethodPatch, ctx)
}

// do send request
func (c *muteHttpClient) do(method string, ctx context.Context) (muteHttpResponse, error) {
	startTime := time.Now().UnixMilli()
	var err error
	var result muteHttpResponse
	c.method = method
	c.request = c.request.WithContext(ctx)
	c.request.URL, err = url.ParseRequestURI(c.url)
	c.request.Method = method
	if err != nil {
		goto RESULT
	}
	result.response, err = c.client.Do(c.request)
	if err != nil {
		goto RESULT
	}
	result.body, err = ioutil.ReadAll(result.response.Body)
	if err != nil {
		goto RESULT
	}
	defer result.response.Body.Close()
	if c.mustCode > 0 && result.response.StatusCode != c.mustCode {
		err = codeMustErr
		goto RESULT
	}
	c.useTime = time.Now().UnixMilli() - startTime
RESULT:
	result.client = *c
	return result, err
}
