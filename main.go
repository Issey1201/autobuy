package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Issey1201/pkg/autobuy"
	"github.com/Issey1201/pkg/notify"
)

// consumerという関数名で良いのか、、、
func consumer(t autobuy.TargetSite, ch chan bool) {
	for {
		if result := autobuy.Check(t); result == false {
			ch <- result
			time.Sleep(1 * time.Minute)
		} else {
			ch <- result
			break
		}
	}
	// closeする箇所がよくわからない
	close(ch)
}

func main() {
	ark := autobuy.NewArk()
	stock := make(chan bool)

	var err error
	go consumer(ark, stock)
	for {
		select {
		case result := <-stock:
			if result == false {
				ark.Tracer.Trace("在庫切れなう")
			} else {
				ark.Tracer.Trace("在庫あったぜ！")
				if err = ark.Run(); err != nil {
					log.Fatalln("Failed to run")
				}
				if err = notify.Notificator(); err != nil {
					fmt.Println("Failed to Notificator")
				}
				os.Exit(1)
			}
		default:
			break
		}
	}
}
