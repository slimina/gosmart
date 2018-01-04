# 发送邮件工具类

## 配置说明
```go
type Config struct { //配置类
	Host      string //邮件服务器地址，如smtp.qq.com
	Port      int  //发送邮件服务器端口25
	Username  string //用户名
	Password  string //密码
	From      string //发送人邮件地址
	FromAlias string //发送人别名
	AuthType  AuthType //邮件服务器认证类型： 枚举：AUTH_TYPE_PLAIN 、AUTH_TYPE_LOGIN AUTH_TYPE_CRAMMD5、默认是AUTH_TYPE_PLAIN
	TplPath   string //邮件模板所在路径
}
```
## 接口说明
```go
type MailService interface {
	SendText(subject string, body string, to ...string) error //发送文本邮件
	SendHtml(subject string, body string, to ...string) error //发送html邮件
	SendTpl(subject string, tplfile string, data map[string]interface{}, to ...string) error //发送模板邮件
	UpdateConfig(Config) //更新邮件配置
}
```
## 使用实例
```go
//初始化配置
cfg := Config{Host: "smtp.qq.com", Username: "286190321@qq.com", Password: "****", Port: 25, From: "286190321@qq.com", FromAlias: "测试例子", AuthType: AUTH_TYPE_PLAIN}
maiService := New(cfg) //实例化发送服务
err := maiService.SendHtml("测试主题", "测试内容<br/>11111", "tianwei7518@qq.com")
if err != nil {
	fmt.Println("test发送失败:" + err.Error())
} else {
	fmt.Println("test发送成功")
}
```