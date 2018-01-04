package mailer

type AuthType int

const (
	AUTH_TYPE_PLAIN   AuthType = iota //0
	AUTH_TYPE_LOGIN                   //1
	AUTH_TYPE_CRAMMD5                 //2
)

type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	FromAlias string
	AuthType  AuthType // 枚举：AUTH_TYPE_PLAIN 、AUTH_TYPE_LOGIN AUTH_TYPE_CRAMMD5、默认是AUTH_TYPE_PLAIN
	TplPath   string
	//SSL       bool
	//Timeout time.Duration //设置超时时间
}

func DefaultConfig() Config {
	return Config{}
}

func (conf Config) IsValid() bool {
	if conf.Host == "" {
		return false
	}
	if conf.Port < 1 || conf.Port > 65535 {
		return false
	}
	if conf.Username == "" || conf.Password == "" {
		return false
	}
	return true
}
