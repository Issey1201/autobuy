package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Issey1201/pkg/autobuy"
	"github.com/Issey1201/pkg/notify"
)

func main() {
	flag.Parse()
	targets := flag.Args()
	// とりあえずarkだけ対応, 将来的に以下みたいなものを想定
	// go run main.go ark amazon newEgg
	ark := autobuy.NewArk(fmt.Sprintf("./config/%v.toml", targets[0]))

	// targetUrlが複数であっても１つだけ買えればそれで良い→チャネルは１つ？

	stockStatus := make(chan bool)
	var err error
	_ = err
	for _, v := range ark.Config.Url.TargetUrl {
		// せっかくのgoなのでtargetUrlの数だけgoroutine回す
		go autobuy.CheckStock(ark, v, stockStatus)
	}

	for v := range stockStatus {
		switch v {
		case true:
			fmt.Println("在庫ある")
			if err = ark.Run(""); err != nil {
				log.Fatalln("Failed to run")
			}
			if err = notify.Notificator(); err != nil {
				fmt.Println("Failed to Notificator")
			}
			os.Exit(1)
		case false:
			fmt.Println("在庫なし")
		}
	}
}
