package netservice

import (
	"log"
	"net/http"
	"testing"
)

type Handler struct {
}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("TTTTTTTT\n"))
}

var handler = &Handler{}

//http://www.2cto.com/kf/201701/587954.html
//http://blog.csdn.net/hherima/article/details/52469674
//handshake_server.go hs.cert, err = c.config.getCertificate(hs.clientHelloInfo())
func Test_Https(t *testing.T) {
	httpsServer := NewHttpsServer()
	httpsServer.Init()
	//	cert, _ := httpsServer.TLSConfig.GetCertificate(nil)
	//	fmt.Println(cert)
	//err := httpsServer.RunHttpsService(":3344", "../config/ca/server.crt", "../config/ca/server.key", handler)
	err := httpsServer.RunHttpsService(":3344", "", "", handler)
	log.Println(err)
}
