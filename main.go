package main

import (
	"autobuy/pkg/autobuy"
	"autobuy/pkg/notify"
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
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

// 流れ
// 1.consumer()で在庫をチェックする
// 2-a.在庫がなければ1分休憩、その後また在庫があるか確認
// 2-b.在庫があればchannelを閉じてRunで購入スクリプト実行
// 3-b.goroutine抜け出してスクリプト終了(通知したい)
// ゆくゆくは、ark以外のサイトも対象とし、ターゲットURLを複数にしていきたい
func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		os.Exit(1)
	}
	ark := autobuy.NewArk()
	user := map[string]string{
		"mailAddress": cfg.Section("ark").Key("user_email").String(),
		"password":    cfg.Section("ark").Key("user_password").String(),
		"name":        cfg.Section("ark").Key("user_name").String(),
		"nameKana":    cfg.Section("ark").Key("user_name_kana").String(),
		"zipcode1":    cfg.Section("ark").Key("user_zipcode1").String(),
		"zipcode2":    cfg.Section("ark").Key("user_zipcode2").String(),
		"pref":        cfg.Section("ark").Key("user_pref").String(),
		"city":        cfg.Section("ark").Key("user_city").String(),
		"street":      cfg.Section("ark").Key("user_street").String(),
		"building":    cfg.Section("ark").Key("user_building").String(),
		"phone":       cfg.Section("ark").Key("user_phone").String(),
		"email":       cfg.Section("ark").Key("user_email").String(),
	}
	stock := make(chan bool)
	go consumer(ark, stock)
	for {
		select {
		case result := <-stock:
			if result == false {
				ark.Tracer.Trace("在庫切れなう")
			} else {
				ark.Tracer.Trace("在庫あったぜ！")
				ark.Run(user)
				notify.Notificator()
				os.Exit(1)
			}
		default:
			break
		}
	}
}
