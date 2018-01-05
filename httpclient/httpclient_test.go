package httpclient

import (
	"fmt"
	"testing"
)

func Test_post_1(t *testing.T) {
	client := New(Config{})
	code, status, body, err := client.HttpGet("https://www.baidu.com")
	if err != nil && code != 200 {
		t.Error("请求失败:" + err.Error())
	} else {
		fmt.Println(code)
		fmt.Println(status)
		fmt.Println(body)
		t.Log("请求成功")
	}
}
