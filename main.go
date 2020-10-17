package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Issey1201/pkg/autobuy"
	"github.com/Issey1201/pkg/notify"
)

func main() {
	flag.Parse()
	targets := flag.Args()
	// とりあえずarkだけ対応, 将来的に以下みたいなものを想定
	// go run main.go ark amazon newEgg
	ark := autobuy.NewArk(fmt.Sprintf("./config/%v.toml", targets[0]))

	ch := make(chan autobuy.CheckResponse, len(ark.Config.Url.TargetUrl))
	var err error

	for _, v := range ark.Config.Url.TargetUrl {
		go autobuy.CheckStock(ark, v, ch)
	}

	for v := range ch {
		switch v.StockStatus {
		case true:
			fmt.Println(time.Now(), " ", v.Url, ": 在庫ある")
			if err = ark.Run(v.Url); err != nil {
				log.Fatalln("Failed to run")
			}
			if err = notify.Notificator(); err != nil {
				fmt.Println("Failed to Notificator")
			}
			close(ch)
			break
		case false:
			fmt.Println(time.Now(), " ", v.Url, ": 在庫なし")
		}
	}
}
