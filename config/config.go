package config

import (
	"encoding/json"
	"io/ioutil"
)

var Config Configure

type Configure struct {
	Addr           string `json:"addr"`
	VerifyInterval int    `json:"verify_interval"`
	ValidateUrl    string `json:"validate_url"`
	Qingting       struct {
		Api           string `json:"api"`
		OrderId       string `json:"order_id"`
		Format        string `json:"format"`
		LineSeparator string `json:"line_separator"`
		CanRepeat     string `json:"can_repeat"`
		Num           int    `json:"num"`
		UserToken     string `json:"user_token"`
	} `json:"qing_ting"`
	KuaiDaili struct {
		UrlTpl string `json:"url_tpl"`
	} `json:"kuai_dai_li"`
}

func InitConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}
	return nil
}
