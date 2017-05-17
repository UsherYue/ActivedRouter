package netservice

import (
	"net/http"
	_ "net/http/pprof"
)

func ListenAndServePProf(addr string, handler http.Handler) {
	http.ListenAndServe(addr, handler)
}
