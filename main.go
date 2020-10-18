package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Issey1201/pkg/autobuy"
	"github.com/Issey1201/pkg/notify"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		log.Fatalf("failed to open env: %v\n", err)
	}
}

func main() {
	flag.Parse()
	targets := flag.Args()
	// とりあえずarkだけ対応, 将来的に以下みたいなものを想定
	// go run main.go ark amazon newEgg
	ark := autobuy.NewArk(fmt.Sprintf("./config/%v.toml", targets[0]))
	//yodobashi := autobuy.NewYodobashi(fmt.Sprintf("./config/%v.toml", targets[0]))

	ch := make(chan autobuy.CheckResponse, len(ark.Config.Url.TargetUrl))
	//ch := make(chan autobuy.CheckResponse, len(yodobashi.Config.Url.TargetUrl))
	done := make(chan struct{})
	var err error

	for _, v := range ark.Config.Url.TargetUrl {
		go autobuy.Check(ark, v, ch, done)
	}
	//for _, v := range yodobashi.Config.Url.TargetUrl {
	//	go autobuy.Check(yodobashi, v, ch, done)
	//}

Loop:
	for v := range ch {
		switch v.StockStatus {
		case true:
			fmt.Println(time.Now().String(), " ", v.Url, ": 在庫ある")
			close(done)
			if err = ark.Run(v.Url); err != nil {
				log.Fatalln("Failed to run")
			}
			//if err = yodobashi.Run(v.Url); err != nil {
			//	log.Fatalln("Failed to run")
			//}
			if err = notify.Notificator(); err != nil {
				fmt.Println("Failed to Notificator")
			}
			close(ch)
			break Loop
		case false:
			fmt.Println(time.Now().String(), " ", v.Url, ": 在庫なし")
		}
	}
}
