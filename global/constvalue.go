package global

//全局数据
//配置文件相关
const (
	ServerConfig    = "config/server.ini"
	HookConfig      = "config/hook.json"
	ClientConfig    = "config/client.json"
	HttpProxyConfig = "config/http_proxy.json"
)

//运行模式相关
const (
	ServerMode       = "server" //服务器模式
	ClientMode       = "client" //客户端模式
	ProxyMode        = "proxy"
	ReserveProxyMode = "reserveproxy" //反向代理模式
	MixMode          = "mix"          //server+proxy模式
	HttpRouterMode   = "httprouter"   //http路由模式
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
		Tencent QQ:4223665
		ActivedRouter Desceiption:
			mode运行模式可以是 Server、Client、Reserveproxy、Mix,分别代表路由服务器模式、客户端模式、反向代理服务器模式、服务器模式&反向代理模式 。
			服务器模式下加载server.ini。
			客户端模式下记加载client.ini读取相关配置信息。	
		The commands are:
			ActiveRouter --runmode=Server/Client/ReserveProxy/Mix
			Client 模式运行在客户端的代理程序
			Server 启动监控服务 
			ReserveProxy 启动反向代理服务
			Mix  启动反向dialing服务并启动监控服务
		The Help:
		    ActiveRouter --help or  -h or -help
		`
	UsageRunmodeTemplate = `{{.Msg}}`
)
