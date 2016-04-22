package global

//全局数据
//配置文件相关
var (
	ServerConfig    = "config/server.ini"
	HookConfig      = "config/hook.json"
	ClientConfig    = "config/client.json"
	HttpProxyConfig = "config/http_proxy.json"
)

//运行模式相关
var (
	ServerMode     = "server"     //服务器模式
	ClientMode     = "client"     //客户端模式
	ProxyMode      = "proxy"      //反向代理模式
	MixMode        = "mix"        //server+proxy模式
	HttpRouterMode = "httprouter" //http路由模式
)

//负载方法

var (
	Alived = "alived"
	Random = "random"
)

// 服务器权重 1-9  权重越高 那么服务器处理能力越高 在服务的时候优先被使用
var (
	HOSTWEIGHT_1 = 1
	HOSTWEIGHT_2 = 2
	HOSTWEIGHT_3 = 3
	HOSTWEIGHT_4 = 4
	HOSTWEIGHT_5 = 5
	HOSTWEIGHT_6 = 6
	HOSTWEIGHT_7 = 7
	HOSTWEIGHT_8 = 8
	HOSTWEIGHT_9 = 9
)
