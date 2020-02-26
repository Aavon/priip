package server

import (
	"net/http"

	"github.com/Aavon/priip/config"
	"github.com/Aavon/priip/ip"
)

func Get(rw http.ResponseWriter, req *http.Request) {
	bestIp := ip.LocalIps.PeekMin()
	if bestIp == nil {
		rw.Write([]byte("no available ip"))
		return
	}
	pv := bestIp.Value.(ip.IPPort)
	rw.Write([]byte(pv.String()))
	return
}

func StratServer() {
	http.HandleFunc("/get", Get)
	http.ListenAndServe(config.Config.Addr, nil)
}
