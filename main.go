package main

import (
	"flag"
	"log"

	"github.com/Aavon/priip/config"
	"github.com/Aavon/priip/ip"
	"github.com/Aavon/priip/server"
)

var (
	CONFIG_FILE string
)

func init() {
	flag.StringVar(&CONFIG_FILE, "c", "./config/priip.json", "config file path")
}

func main() {
	flag.Parse()
	if err := config.InitConfig(CONFIG_FILE); err != nil {
		log.Fatal(err)
	}
	ip.InitLocalIps()
	//qingting := ip.NewQingtingIP(
	//	config.Config.Qingting.Api,
	//	config.Config.Qingting.OrderId,
	//	config.Config.Qingting.Format,
	//	config.Config.Qingting.LineSeparator,
	//	config.Config.Qingting.CanRepeat,
	//	config.Config.Qingting.UserToken,
	//	config.Config.Qingting.Num,
	//)
	kuaidaili := ip.NewKuaidailiFree(config.Config.KuaiDaili.UrlTpl)
	verifier := ip.NewVerifer(
		[]ip.IPPuller{kuaidaili},
		config.Config.VerifyInterval,
		config.Config.ValidateUrl,
	)
	go verifier.Start()
	server.StratServer()
}
