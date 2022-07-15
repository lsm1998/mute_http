package mute

import (
	"context"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	response, err := New("https://www.baidu.com").Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.Code())
	fmt.Println(response.Curl())
	fmt.Println(response.UseTime())
	fmt.Println(string(response.GetBody()))
}
