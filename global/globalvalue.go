package global

//全局数据
//配置文件相关
var (
	ServerConfig = "config/server.ini"
	HookConfig   = "config/hook.json"
	ClientConfig = "config/client.ini"
	ProxyConfig  = "config/proxy.json"
)

//运行模式相关
var (
	ServerMode     = "server"     //服务器模式
	ClientMode     = "client"     //客户端模式
	ProxyMode      = "proxy"      //反向代理模式
	MixMode        = "mix"        //server+proxy模式
	HttpRouterMode = "httprouter" //http路由模式
)
