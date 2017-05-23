package global

//config file
const (
	ServerJsonConfig = "config/server.json"
	HookConfig       = "config/hook.json"
	ClientConfig     = "config/client.json"
	HttpProxyConfig  = "config/http_proxy.json"
	HttpsProxyConfig = "config/https_proxy.json"
	CertificateData  = "config/crtdata"
)

const (
	//http statistics interval 5min
	Http_Statistics_Interval = 60
	//defaule pprof's http service  address
	HTTP_PPROF_Default_Addr = ":6060"
	DefaultHttpAddr         = "127.0.0.1:80"
	DefaultHttsAddr         = "127.0.0.1:443"
)

//certificate
const (
	DefaultCertificate = "server.crt"
	DefaultKey         = "server.key"
)

//run mode
const (
	ServerMode       = "server"       //server mode
	ClientMode       = "client"       //client mode
	ReverseProxyMode = "reverseproxy" //reverseproxy mode
	InitMode         = "init"         //初始化
)

//load balance method
const (
	Alived = "alived"
	Random = "random"
)

const (
	SwitchOn  = "on"
	SwitchOff = "off"
)

//Server weight,The higher the weight, the easier it is to be scheduled
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
	DescTemplate  = "\n\033[35;1mActivedRouter Description:\033[0m\n\033[1mActiveRouter is a load balancing software and provides a web dashboard for management.\nAuthor:usher.yue\nEmail:ushe.yue@gmail.com\033[0m\033[35;1m\n\nUSAGE:\033[0m\n\t\033[1mActiveRouter --runmode= [arguments]\033[0m\n  \033[35;1m\nThe commands are:\033[0m\n"
	UsageTemplate = `
	ActiveRouter --runmode=Server  Running In Server Mode 
	ActiveRouter --runmode=Client  Running In Client Mode
	ActiveRouter --runmode=ReverseProxy   Running In ReserveProxy Mode
`
	TheHelpTemplate   = "\033[35;1mThe Help:\033[0m\n\033[1mActiveRouter --help or  -h or -help\033[0m"
	UsageRunmodeError = "runmode parameters error ,please reference  ActiveRouter --runmode=Client/ReverseProxy/Server"
)
