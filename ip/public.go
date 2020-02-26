package ip

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type IPPort struct {
	Ip   string
	Port string
}

func (i *IPPort) String() string {
	return fmt.Sprintf("%s:%s", i.Ip, i.Port)
}

func (i *IPPort) URL() string {
	return fmt.Sprintf("http://%s:%s", i.Ip, i.Port)
}

type IPPuller interface {
	// 获取完整IP列表
	Pull() ([]IPPort, error)
}

type QingtingIP struct {
	// https://proxyapi.horocn.com/api/v2/proxies
	Api string
	// 订单ID
	OrderId string
	Num     int
	// text json
	Format string
	// win unix
	LineSeparator string
	// yes no
	CanRepeat string
	UserToken string
	uri       string
	c         *http.Client
}

type QingtingIPResp struct {
	// 0
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type QingtingIPInfo struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	CountryCn  string `json:"country_cn"`
	ProvinceCn string `json:"province_cn"`
	CityCn     string `json:"city_cn"`
}

func NewQingtingIP(api, orderId, format, lineSeparator, canRepeat, userToken string, num int) *QingtingIP {
	q := &QingtingIP{
		Api:           api,
		OrderId:       orderId,
		Format:        format,
		LineSeparator: lineSeparator,
		CanRepeat:     canRepeat,
		UserToken:     userToken,
		Num:           num,
	}
	q.uri = fmt.Sprintf("%s?order_id=%s&num=%d&format=%s&line_separator=%s&can_repeat=%s&user_token=%s",
		q.Api,
		q.OrderId,
		q.Num,
		q.Format,
		q.LineSeparator,
		q.CanRepeat,
		q.UserToken,
	)
	q.c = &http.Client{
		Timeout: 10 * time.Second,
	}
	return q
}

func (q *QingtingIP) Pull() ([]IPPort, error) {
	log.Println(q.uri)
	r, err := q.c.Get(q.uri)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", data)
	resp := QingtingIPResp{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("qingting %d", context.DeadlineExceeded)
	}
	ipList := []QingtingIPInfo{}
	err = json.Unmarshal(resp.Data, &ipList)
	if err != nil {
		return nil, err
	}
	result := make([]IPPort, 0, len(resp.Data))
	for _, p := range ipList {
		result = append(result, IPPort{
			Ip:   p.Host,
			Port: p.Port,
		})
	}
	return result, nil
}
