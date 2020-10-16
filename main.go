package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/Issey1201/pkg/autobuy"
	"github.com/Issey1201/pkg/notify"
)

func main() {
	flag.Parse()
	targets := flag.Args()
	// とりあえずarkだけ対応, 将来的に以下みたいなものを想定
	// go run main.go ark amazon newEgg
	ark := autobuy.NewArk(fmt.Sprintf("./config/%v.toml", targets[0]))

	wg := new(sync.WaitGroup)
	ch := make(chan autobuy.CheckResponse, len(ark.Config.Url.TargetUrl))
	var err error

	for _, v := range ark.Config.Url.TargetUrl {
		// せっかくのgoなのでtargetUrlの数だけgoroutine回す
		wg.Add(1)
		go autobuy.CheckStock(ark, v, ch, wg)
	}

	for v := range ch {
		switch v.StockStatus {
		case true:
			fmt.Println(v.Url, ": 在庫ある")
			if err = ark.Run(v.Url); err != nil {
				log.Fatalln("Failed to run")
			}
			if err = notify.Notificator(); err != nil {
				fmt.Println("Failed to Notificator")
			}
			break
		case false:
			fmt.Println(v.Url, ": 在庫なし")
		}
	}
	//wg.Wait()
	fmt.Println("ぬけたた")
	//for {
	//	select {
	//	case res := <- ch:
	//		switch res.StockStatus {
	//		case true:
	//			wg.Wait()
	//			fmt.Println(res.Url, ": 在庫ある")
	//			if err = ark.Run(res.Url); err != nil {
	//				log.Fatalln("Failed to run")
	//			}
	//			if err = notify.Notificator(); err != nil {
	//				fmt.Println("Failed to Notificator")
	//			}
	//			os.Exit(1)
	//		case false:
	//			fmt.Println(res.Url, ": 在庫なし")
	//		}
	//	}
	//}
}
