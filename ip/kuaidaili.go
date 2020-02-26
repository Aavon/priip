package ip

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 最大页码
const MAX_PAGE = 5

type KuaidailiFree struct {
	// https://www.kuaidaili.com/free/inha/%d
	UrlTemplate string
}

func NewKuaidailiFree(urltpl string) *KuaidailiFree {
	return &KuaidailiFree{
		UrlTemplate: urltpl,
	}
}

func (k *KuaidailiFree) Pull() ([]IPPort, error) {
	result := make([]IPPort, 0)
	for i := 1; i <= MAX_PAGE; i++ {
		//User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0
		req, err := http.NewRequest("GET", fmt.Sprintf(k.UrlTemplate, i), nil)
		if err != nil {
			return result, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return result, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return result, err
		}
		doc.Find("#content tbody tr").Each(func(i int, selection *goquery.Selection) {
			ipCol := selection.Find(" [data-title=IP]")
			portCol := selection.Find(" [data-title=PORT]")
			item := IPPort{}
			item.Ip = ipCol.Text()
			item.Port = portCol.Text()
			result = append(result, item)
		})
		time.Sleep(time.Second)
	}
	return result, nil
}
