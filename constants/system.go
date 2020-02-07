package constants

const (
	// ------ 系统 ------
	// cookie
	DaoSuanSystemCookie = "DAOSUAN_SYSTEM_COOKIE"
	// session过期时间 10小时
	DaoSuanSessionExpires = 864000
	// cookie过期时间  7天
	DaoSuanCookieExpires = 3600 * 24 * 7
	// 登陆态session名称
	DaoSuanSessionName = "DAOSUAN_AUTHENTICATION_KEY"

	// ------ 资源 ------
	StorageTokenTime = 300


	// 任务队列长度
	QueueTaskLength = 100
)
