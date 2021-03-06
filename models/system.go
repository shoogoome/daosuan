package models

import (
	wechat2 "github.com/shoogoome/gowechat"
	"golang.org/x/oauth2"
)

type SystemConfiguration struct {
	Mysql         mysql         `json:"mysql" yaml:"mysql"`                 // 数据库配置
	Server        server        `json:"server" yaml:"server"`               // 系统配置
	Dijan         dijan         `json:"dijan" yaml:"dijan"`                 // 缓存配置
	QiNiu         qiniu         `json:"qi_niu" yaml:"qi_niu"`               // 七牛配置
	Nsq           nsq           `json:"nsq" yaml:"nsq"`                     // nsq配置
	Elasticsearch elasticsearch `json:"elasticsearch" yaml:"elasticsearch"` // elasticsearch配置
	Oauth         oauth         `json:"oauth" yaml:"oauth"`                 // Oauth
}

type server struct {
	Salt                       string `json:"salt" yaml:"salt"`                                                     // 盐
	TokenBucketCapacity        int64  `json:"token_bucket_capacity" yaml:"token_bucket_capacity"`                   // 令牌桶允许最大大小，即允许瞬间爆发请求
	TokenBucketOutputPerSecond int    `json:"token_bucket_output_per_second" yaml:"token_bucket_output_per_second"` // 令牌桶每秒产出，qps
	TaskQueueLength            int    `json:"task_queue_length" yaml:"task_queue_length"`                           // 任务队列长度
	Mail                       mail   `json:"mail" yaml:"mail"`                                                     // 邮箱服务配置
}

type mysql struct {
	DB       string `json:"db" yaml:"db"`             // db名
	Host     string `json:"host" yaml:"host"`         // 主机名
	Port     string `json:"port" yaml:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
}

type dijan struct {
	Host       string `json:"host" yaml:"host"`               // 主机名
	Port       int    `json:"port" yaml:"port"`               // 端口
	Node       int    `json:"node" yaml:"node"`               // 总节点数
	PoolNumber int    `json:"pool_number" yaml:"pool_number"` // 连接池数量
}

type qiniu struct {
	Bucket    string `json:"bucket" yaml:"bucket"`         // 空间名
	Expires   uint64 `json:"expires" yaml:"expires"`       // 过期时间
	AccessKey string `json:"access_key" yaml:"access_key"` // 密钥
	SecretKey string `json:"secret_key" yaml:"secret_key"` // 密钥
}

type nsq struct {
	Host string `json:"host" yaml:"host"` // host
	Port int    `json:"port" yaml:"port"` // 端口
}

type elasticsearch struct {
	Host string `json:"host" yaml:"host"` // host
	Port int    `json:"port" yaml:"port"` // port
}

type oauth struct {
	GitHub github `json:"github" yaml:"github"` // github验证
	WeChat wechat                               // 微信验证
}

type github struct {
	ClientId     string `json:"client_id" yaml:"client_id"`         // client id
	ClientSecret string `json:"client_secret" yaml:"client_secret"` // client secret
	CookieDomain string `json:"cookie_domain" yaml:"cookie_domain"` // Cookie domain
	RedirectUrl  string `json:"redirect_url" yaml:"redirect_url"`   // 回调地址
	SuccessUrl   string `json:"success_url" yaml:"success_url"`     // 验证成功回调地址
	ErrorUrl     string `json:"error_url" yaml:"error_url"`         // 验证失败回调地址
	Oauth2Config oauth2.Config                                      // oauth配置
}

type wechat struct {
	OauthClient wechat2.WeCharClient
}

type mail struct {
	SmtpHost string `json:"smtp_host" yaml:"smtp_host"` // 邮箱服务地址
	Username string `json:"username" yaml:"username"`   // 账号
	Password string `json:"password" yaml:"password"`   // 密码
}
