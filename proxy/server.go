package proxy

import (
	"sync"

	"github.com/Aavon/priip/proxy/config"
	"github.com/Aavon/priip/proxy/mitm"
)

func StartServer() {
	conf := config.Cfg{}
	// init tls config
	tlsConfig := config.NewTlsConfig("gomitmproxy-ca-pk.pem", "gomitmproxy-ca-cert.pem", "", "")
	// start mitm proxy
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mitm.Gomitmproxy(&conf, tlsConfig, wg)
	wg.Wait()
}
