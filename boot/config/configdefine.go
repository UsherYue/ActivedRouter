package config

//client config mapping struct
type ClientConfigData struct {
	Domain           string   `json:"domain"`
	ClusterName      string   `json:"cluster"`
	RouterServerList []string `json:"router_list"`
}

//server config mapping struct
type ServerConfigData struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	ServerMode string `json:"srvmode"`
	HttpHost   string `json:"httphost"`
	HttpPort   string `json:"httpport"`
}
