package autobuy

import (
	"log"
	"os"
	"testing"

	"github.com/go-ini/ini"
)

func TestArk_Run(t *testing.T) {

	// 相対パスじゃなくて絶対パスにしたい
	cfg, err := ini.Load("./../../config.ini")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		os.Exit(1)
	}
	user := map[string]string{
		"mailAddress": cfg.Section("ark").Key("user_email").String(),
		"password": cfg.Section("ark").Key("user_password").String(),
	}

	// ログイン情報とURLは引数とするべきか？それとも構造体に格納すべきか？
	// arkに関する固定情報は構造体、user情報などarkでも情報が状況により変わってくるのは引数が良い？
	ark := NewArk()
	if err := ark.Run(user); err != nil {
		t.Errorf("errorを返すべきでは無い: %v", err)
	}
}