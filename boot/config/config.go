package config

import (
	. "ActivedRouter/global"
	"ActivedRouter/hook"
	"ActivedRouter/netservice"
)

//parse config
func ParseConfigfile() {
	switch RunMode {
	case ServerMode:
		{
			//server config
			netservice.LoadServerJsonConfig(ServerJsonConfig)
			//hook script
			hook.ParseHookScript(HookConfig)
		}
	case ClientMode:
		{
			//client mode
			netservice.LoadClientConfig(ClientConfig)
		}
	case ReverseProxyMode:
		{
			//server config
			netservice.LoadServerJsonConfig(ServerJsonConfig)
			//certificate config
			netservice.DefaultReverseProxy.LoadCertificateConfig(CertificateData)
			//proxy config
			netservice.DefaultReverseProxy.LoadProxyConfig(HttpProxyConfig)
		}
	case InitMode:
		{
			//init app
			//init config file
		}
	}
}
