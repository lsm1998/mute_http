package mute

import (
	"context"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	response, err := New("").Get(context.Background())
	fmt.Println(response, err)
}
