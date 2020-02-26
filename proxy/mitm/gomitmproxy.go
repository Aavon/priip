// This example shows a proxy server that uses go-mitm to man-in-the-middle
// HTTPS connections opened with CONNECT requests

package mitm

import (
	"net/http"
	"sync"
	"time"

	"github.com/Aavon/priip/proxy/config"
	"log"
)

func Gomitmproxy(conf *config.Cfg, tlsConfig *config.TlsConfig, wg *sync.WaitGroup) {
	handler, err := InitConfig(conf, tlsConfig)
	if err != nil {
		log.Fatalf("InitConfig error: %s", err)
	}

	server := &http.Server{
		Addr:         ":" + *conf.Port,
		Handler:      handler,
		ReadTimeout:  1 * time.Hour,
		WriteTimeout: 1 * time.Hour,
	}

	go func() {
		log.Printf("Gomitmproxy Listening On: %s", *conf.Port)
		if *conf.Tls {
			log.Println("Listen And Serve HTTP TLS")
			err = server.ListenAndServeTLS("gomitmproxy-ca-cert.pem", "gomitmproxy-ca-pk.pem")
		} else {
			log.Println("Listen And Serve HTTP")
			err = server.ListenAndServe()
		}
		if err != nil {
			log.Fatalf("Unable To Start HTTP proxy: %s", err)
		}

		wg.Done()
		log.Printf("Gomitmproxy Stop!!!!")
	}()

	return
}
