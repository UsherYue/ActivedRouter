package global

//全局数据
//配置文件相关
const (
	ServerConfig     = "config/server.ini"
	ServerJsonConfig = "config/server.json"
	HookConfig       = "config/hook.json"
	ClientConfig     = "config/client.json"
	HttpProxyConfig  = "config/http_proxy.json"
)

//运行模式相关
const (
	ServerMode       = "server"       //服务器模式
	ClientMode       = "client"       //客户端模式
	ReserveProxyMode = "reserveproxy" //反向代理模式
)

//负载方法
const (
	Alived = "alived"
	Random = "random"
)

// 服务器权重 1-9  权重越高 那么服务器处理能力越高 在服务的时候优先被使用
const (
	HOSTWEIGHT_1 = iota + 1
	HOSTWEIGHT_2
	HOSTWEIGHT_3
	HOSTWEIGHT_4
	HOSTWEIGHT_5
	HOSTWEIGHT_6
	HOSTWEIGHT_7
	HOSTWEIGHT_8
	HOSTWEIGHT_9
)

const (
	UsageTemplate = `
		ActiveRouter 是一个简单的基于路由分发和反向代理的负载均衡监控服务,并且提供一个方便管理的仪表盘用于快速配置、挂载、卸载.
		Author:usher.yue
		Email:ushe.yue@gmail.com
		ActivedRouter Desceiption:
			mode运行模式可以是 Server、Client、Reserveproxy,分别代表路由服务器模式、客户端模式、反向代理服务器模式。
		The commands are:
			ActiveRouter --runmode=Server/Client/ReserveProxy
			Client 启动客户端模式,一般运行在需要监控、代理的服务器上
			Server 启动服务器模式
			ReserveProxy 启动反向代理服务,可选择性的开启关闭服务器模式
		The Help:
		    ActiveRouter --help or  -h or -help
		`
	UsageRunmodeError = `"runmode参数错误,参考 ActiveRouter --runmode=Client/Reserveproxy/Server"`
)
