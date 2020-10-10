package autobuy

import (
	"autobuy/pkg/trace"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
	"github.com/sclevine/agouti"
)

type Ark struct {
	url        string
	targetUrl  string
	inputMail  string
	inputPw    string
	login      string
	checkPoint string
	checkWord  string
	Tracer     trace.Tracer
}

func NewArk() *Ark {
	// 相対パスじゃなくて絶対パスにしたい
	cfg, err := ini.Load("./../../config.ini")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		os.Exit(1)
	}
	return &Ark{
		url: cfg.Section("ark").Key("url").String(),
		targetUrl: cfg.Section("ark").Key("target_url").String(),
		inputMail: cfg.Section("ark").Key("input_xpath_email").String(),
		inputPw: cfg.Section("ark").Key("input_xpath_password").String(),
		login: cfg.Section("ark").Key("input_xpath_login").String(),
		checkPoint: cfg.Section("ark").Key("selector_stock").String(),
		checkWord: cfg.Section("ark").Key("out_of_stock_word").String(),
		Tracer : trace.New(os.Stdout),
	}
}

var ChooseArk Ark

func (t *Ark) Run (user map[string]string) (err error){
	// ブラウザ：chromeを指定して起動
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
		return err
	}
	if driver != nil {
		defer driver.Stop()
	}

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page: %v", err)
		return err
	}

	if err := page.Navigate(t.url); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
		return err
	}
	sleep()

	if err := page.ClearCookies(); err != nil {
		log.Fatalf("Failed to clear cookies: %v", err)
		return err
	}
	sleep()

	email := page.FindByID(t.inputMail)
	password := page.FindByXPath(t.inputPw)
	email.Fill(user["mailAddress"])
	password.Fill(user["password"])
	if err := page.FindByXPath(t.login).Submit();
		err != nil {
			log.Fatalf("Failed to login: %v", err)
			return err
	}
	sleep()

	if err := page.Navigate(t.targetUrl); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
		return err
	}
	sleep()

	if result := t.Check(); result == false {
		return errors.New("在庫ないです")
	}else{
		return errors.New("在庫あります")
	}

	return nil
}

func sleep (){
	time.Sleep(1 * time.Second)
	// 以下のようにするとエラー起きる、、、
	//time.Sleep(s * time.Second)
}