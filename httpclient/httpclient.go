package httpclient

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient interface { //定义接口
	HttpPost(sendurl string, contentType string, requestBody io.Reader) (code int, status string, body string, err error)
	HttpPostForm(sendurl string, data url.Values) (code int, status string, body string, err error)
	HttpGet(url string) (code int, status string, body string, err error)
}

//配置实体
type Config struct {
	ConnectTimeout        time.Duration
	RequestTimeout        time.Duration
	KeepAlive             time.Duration
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	ResponseHeaderTimeout time.Duration
}

//定义实体
type client struct {
	*http.Client
}

func New(cfg Config) HttpClient {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   cfg.ConnectTimeout * time.Second,
				KeepAlive: cfg.KeepAlive * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          cfg.MaxIdleConns,
			IdleConnTimeout:       cfg.IdleConnTimeout * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: time.Second * cfg.ResponseHeaderTimeout,
		},
	}
	return &client{c}
}

func (c *client) HttpPost(sendurl string, contentType string, requestBody io.Reader) (code int, status string, body string, err error) {
	var resp *http.Response
	resp, err = c.Post(sendurl, contentType, requestBody)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	code = resp.StatusCode
	status = resp.Status
	if resp.Body != nil {
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		body = string(data)
	}
	return
}

func (c *client) HttpPostForm(sendurl string, data url.Values) (code int, status string, body string, err error) {
	var resp *http.Response
	resp, err = c.PostForm(sendurl, data)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}

	code = resp.StatusCode
	status = resp.Status

	if resp.Body != nil {
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		body = string(data)
	}
	return
}

func (c *client) HttpGet(url string) (code int, status string, body string, err error) {
	var resp *http.Response
	resp, err = c.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	code = resp.StatusCode
	status = resp.Status
	if resp.Body != nil {
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		body = string(data)
	}
	return
}

func GetRemoteIp(req *http.Request) (ip string, port string) {
	h := req.Header
	ip = h.Get("X-Forwarded-For")
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = h.Get("Proxy-Client-IP")
	}
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = h.Get("WL-Proxy-Client-IP")
	}
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = h.Get("HTTP_CLIENT_IP")
	}
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = h.Get("HTTP_X_FORWARDED_FOR")
	}
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = h.Get("X-Real-IP")
	}
	if len(ip) == 0 || strings.EqualFold("unknown", ip) {
		ip = req.RemoteAddr
	}

	ips := strings.Split(ip, ",")
	if len(ips) == 0 {
		//ip = ip
	} else if len(ips) == 1 {
		ip = ips[0]
	} else {
		for _, v := range ips {
			if len(v) == 0 || strings.EqualFold("unknown", v) {
				continue
			}
			ip = v
		}
	}

	if ip != "" {
		ips_ := strings.Split(ip, ":")
		ip = ips_[0]
		if len(ips_) > 1 {
			port = ips_[1]
		}
	}
	return
}
