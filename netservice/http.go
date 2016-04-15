//对于客户端 服务器以及http服务的网络封装
package netservice

import (
	"ActivedRouter/global"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Http struct {
	Host string
	Port string
}

var templates = `
<html>



</html>
`

//ClientInfos
func ClientInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := global.GHostInfoTable.HostsInfo.GetMemory().GetData()
	bts, _ := json.MarshalIndent(data, "", " ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bts))
}

//index redirect to static
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static", 302)
}

//domain
func TestDomain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//w.Write([]byte("hello,world"))
	bts, _ := json.Marshal(r.Host)
	w.Write(bts)
}

//创建http服务
func NewHttp(host, port string) *Http {
	return &Http{Host: host, Port: port}
}

//run router server
func (seft *Http) Run() {
	log.Printf("开始启动http服务,%s:%s........\n", seft.Host, seft.Port)
	router := httprouter.New()
	router.GET("/domain", TestDomain)
	router.GET("/clientinfos", ClientInfos)
	router.GET("/", Index)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	log.Fatal(http.ListenAndServe(seft.Host+":"+seft.Port, router))
}
