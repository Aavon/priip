package ip

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// for test

const ValidateApi = "http://httpbin.org/ip"

func FilterIP(ips []IPPort) []IPPort {
	filtered := []IPPort{}
	for _, ip := range ips {
		// 得加上http
		proxyUrl, _ := url.Parse(ip.URL())
		c := http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
				Proxy:               http.ProxyURL(proxyUrl),
			},
		}
		r, err := c.Get(ValidateApi)
		if err != nil {
			log.Printf("fiter ip [%s]: %v", ip.String(), err)
			continue
		}
		data, _ := ioutil.ReadAll(r.Body)
		log.Printf("filter ip resp [%s]: %s", ip.String(), data)
		r.Body.Close()
		filtered = append(filtered, ip)
	}
	return filtered
}
