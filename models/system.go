package models

type SystemConfiguration struct {

	// 数据库配置
	Mysql mysql `json:"mysql" yaml:"mysql"`

	// 系统配置
	Server server `json:"server" yaml:"server"`

	// 缓存配置
	Dijan dijan `json:"dijan" yaml:"dijan"`

	// 七牛配置
	QiNiu qiniu `json:"qi_niu" yaml:"qi_niu"`
}


type server struct {

	// 盐
	Salt string `json:"salt" yaml:"salt"`

	// 令牌桶允许最大大小，即允许瞬间爆发请求
	TokenBucketCapacity int64 `json:"token_bucket_capacity" yaml:"token_bucket_capacity"`

	// 令牌桶每秒产出，qps
	TokenBucketOutputPerSecond int `json:"token_bucket_output_per_second" yaml:"token_bucket_output_per_second"`

	// 任务队列长度
	TaskQueueLength int `json:"task_queue_length" yaml:"task_queue_length"`
}

type mysql struct {

	// db名
	DB string `json:"db" yaml:"db"`

	// 主机名
	Host string `json:"host" yaml:"host"`

	// 端口
	Port string `json:"port" yaml:"port"`

	// 用户名
	Username string `json:"username" yaml:"username"`

	// 密码
	Password string `json:"password" yaml:"password"`

}

type dijan struct {

	// 主机名
	Host string `json:"host" yaml:"host"`

	// 端口
	Port int `json:"port" yaml:"port"`

	// 总节点数
	Node int `json:"node" yaml:"node"`

	// 连接池数量
	PoolNumber int `json:"pool_number" yaml:"pool_number"`

}

type qiniu struct {

	// 空间名
	Bucket string `json:"bucket" yaml:"bucket"`

	// 过期时间
	Expires uint64 `json:"expires" yaml:"expires"`

	// 密钥
	AccessKey string `json:"access_key" yaml:"access_key"`

	// 密钥
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}