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
