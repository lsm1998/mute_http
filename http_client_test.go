package mute

import (
	"context"
	"fmt"
	"net/http/cookiejar"
	"testing"
)

func TestNew(t *testing.T) {
	response, err := New("https://www.baidu.com").Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	showResponse(response)
}

func TestMuteHttpClient_SetCookieJar(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	client := New("https://www.baidu.com")
	_, err := client.SetCookieJar(jar).Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Get(context.Background())
	fmt.Println(response.Request().Cookies())
	fmt.Println(response.Curl())
}

func showResponse(response muteHttpResponse) {
	fmt.Println(response.Code())
	fmt.Println(response.Curl())
	fmt.Println(response.UseTime())
	fmt.Println(string(response.GetBody()))
}
