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
			netservice.LoadClientJsonConfig(ClientConfig)
		}
	case ReverseProxyMode:
		{
			//server config
			netservice.LoadServerJsonConfig(ServerJsonConfig)
			//certificate config
			netservice.ProxyHandler.LoadCertificateConfig(CertificateData)
			//proxy config
			netservice.ProxyHandler.LoadProxyConfig(HttpProxyConfig)
		}
	case InitMode:
		{
			//init app
			//init config file
		}
	}
}
