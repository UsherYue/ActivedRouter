package boot

//客户端配置
type ClientConfigData struct {
	Domain           string   `json:"domain"`
	ClusterName      string   `json:"cluster"`
	RouterServerList []string `json:"router_list"`
}

//服务器配置
type ServerConfigData struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	ServerMode string `json:"srvmode"`
	HttpHost   string `json:"httphost"`
	HttpPort   string `json:"httpport"`
}

//负载均衡节点
type LbNode struct {
	Domain  string              `json:"domain"`
	Clients []map[string]string `json:"clients"`
}

//反向代理配置
type ReserveProxyConfigData struct {
	ProxyMethod    string  `json:"proxy_method"`
	HttpProxyAddr  string  `json:"http_proxy_addr"`
	HttpSwitch     string  `json:"http_switch"`
	HttpsSwitch    string  `json:"https_switch"`
	HttpsCrt       string  `json:"https_crt"`
	HttpsKey       string  `json:"https_key"`
	HttpsProxyAddr string  `json:"https_proxy_addr"`
	ReserveProxy   *LbNode `json:"reserve_proxy"`
}
