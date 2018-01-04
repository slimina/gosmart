package mailer

import (
	"testing"
)

func Test_mail_1(t *testing.T) {
	cfg := Config{Host: "smtp.qq.com", Username: "286190321@qq.com", Password: "****", Port: 25, From: "286190321@qq.com", FromAlias: "测试例子", AuthType: AUTH_TYPE_PLAIN}
	maiService := New(cfg)
	err := maiService.SendText("测试主题", "测试内容<br/>11111", "tianwei7518@qq.com")
	if err != nil {
		t.Error("test2发送失败:" + err.Error())
	} else {
		t.Log("test2发送成功")
	}
}
func Test_mail_2(t *testing.T) {
	cfg := Config{Host: "smtp.qq.com", Username: "286190321@qq.com", Password: "****", Port: 25, From: "286190321@qq.com", FromAlias: "测试例子"}
	maiService := New(cfg)
	err := maiService.SendHtml("测试主题", "测试内容<br/>11111", "tianwei7518@qq.com")
	if err != nil {
		t.Error("test3发送失败:" + err.Error())
	} else {
		t.Log("test3发送成功")
	}
}

func Test_mail_3(t *testing.T) {
	cfg := Config{Host: "smtp.qq.com", Username: "286190321@qq.com", Password: "****", Port: 25, From: "286190321@qq.com", FromAlias: "测试例子",
		TplPath: "./tpl"}
	maiService := New(cfg)
	data := make(map[string]interface{})
	data["name"] = "lucy"
	maiService.SendTpl("测试主题", "index.html", data, "tianwei7518@qq.com")
	if err != nil {
		t.Error("test3发送失败:" + err.Error())
	} else {
		t.Log("test3发送成功")
	}
}
