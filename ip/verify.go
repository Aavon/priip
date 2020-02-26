package ip

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const INVALID_IP_SCORE = SCORE(0)

type Verifier struct {
	Publics     []IPPuller
	IntervalSec int
	// http://httpbin.org/ip
	ValidateUrl string
	stop        chan bool
}

func NewVerifer(publics []IPPuller, interval int, validateUrl string) *Verifier {
	v := &Verifier{
		Publics:     publics,
		IntervalSec: interval,
		ValidateUrl: validateUrl,
		stop:        make(chan bool, 1),
	}
	return v
}

func (v *Verifier) Start() {
	v.RefreshIps()
	t := time.NewTicker(time.Second * time.Duration(v.IntervalSec))
	for {
		select {
		case <-t.C:
			v.RefreshIps()
		case <-v.stop:
			return
		}
	}
}

func (v *Verifier) Stop() {
	v.stop <- true
}

func (v *Verifier) RefreshIps() {
	all := []IPPort{}
	for _, p := range v.Publics {
		ips, err := p.Pull()
		if err != nil {
			log.Printf("refresh pull err: %v", err)
			continue
		}
		all = append(all, ips...)
	}
	log.Printf("pull count: %d\n", len(all))
	filtered := v.validate(all)
	log.Printf("validate count: %d\n", len(filtered))
	ts := time.Now().UnixNano()
	for _, f := range filtered {
		// 新检查的, score越小
		score := f.Consume - ts
		LocalIps.AddOrUpdate(f.String(), SCORE(score), f.IPPort)
	}
	// 清理不可用IP
	for {
		last := LocalIps.PeekMax()
		if last == nil {
			break
		}
		if last.Score() < 0 {
			break
		}
		// 不可用标记
		if last.Score() == INVALID_IP_SCORE || last.Value == nil {
			LocalIps.Remove(last.Key())
		}
	}
	log.Printf("verify at %d", ts)
}

type ValidateInfo struct {
	IPPort
	// 耗时
	Consume int64
}

func (v *Verifier) validate(ips []IPPort) []ValidateInfo {
	filtered := []ValidateInfo{}
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
		s0 := time.Now().UnixNano()
		r, err := c.Get(v.ValidateUrl)
		if err != nil {
			log.Printf("fiter ip [%s]: %v", ip.String(), err)
			continue
		}
		data, _ := ioutil.ReadAll(r.Body)
		log.Printf("filter ip resp [%s]: %s", ip.String(), data)
		r.Body.Close()
		s1 := time.Now().UnixNano()
		filtered = append(filtered, ValidateInfo{
			IPPort:  ip,
			Consume: s1 - s0,
		})
	}
	return filtered
}
