package autobuy

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Issey1201/pkg/trace"
	"github.com/go-ini/ini"
	"github.com/sclevine/agouti"
)

type Ark struct {
	url          string
	targetUrl    string
	addresseeUrl string
	inputMail    string
	inputPw      string
	login        string
	stockBtn     string
	checkWord    string
	name         string
	nameKana     string
	zipcode1     string
	zipcode2     string
	pref         string
	city         string
	street       string
	building     string
	phone        string
	userEmail    string
	vUserEmail   string
	shipping     string
	payment      string
	nextPage1    string
	nextPage2    string
	nextPage3    string
	Tracer       trace.Tracer
}

type TargetSite interface {
	Run(map[string]string) error
	getCheckInfo() map[string]string
}

func NewArk() *Ark {
	// 相対パスじゃなくて絶対パスにしたい
	// →main.goでもini.Load()をしていて、パスはmain.goに合わせてる感じ？
	// main.goでは相対パスなのに大して、こちらは絶対パス（？）少なくとも相対パスではない
	// 最初読み込んだら２回目以降の読み込みは、１回目と同じ読み込み方で良いということ？
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		os.Exit(1)
	}
	return &Ark{
		url:          cfg.Section("ark").Key("base_url").String(),
		targetUrl:    cfg.Section("ark").Key("target_url").String(),
		addresseeUrl: cfg.Section("ark").Key("addressee_url").String(),
		inputMail:    cfg.Section("ark").Key("input_xpath_email").String(),
		inputPw:      cfg.Section("ark").Key("input_xpath_password").String(),
		login:        cfg.Section("ark").Key("input_xpath_login").String(),
		stockBtn:     cfg.Section("ark").Key("selector_stock").String(),
		checkWord:    cfg.Section("ark").Key("in_stock_word").String(),
		name:         cfg.Section("ark").Key("input_xpath_name").String(),
		nameKana:     cfg.Section("ark").Key("input_xpath_name_kana").String(),
		zipcode1:     cfg.Section("ark").Key("input_xpath_zipcode1").String(),
		zipcode2:     cfg.Section("ark").Key("input_xpath_zipcode2").String(),
		pref:         cfg.Section("ark").Key("input_xpath_pref").String(),
		city:         cfg.Section("ark").Key("input_xpath_city").String(),
		street:       cfg.Section("ark").Key("input_xpath_street").String(),
		building:     cfg.Section("ark").Key("input_xpath_building").String(),
		phone:        cfg.Section("ark").Key("input_xpath_phone").String(),
		userEmail:    cfg.Section("ark").Key("input_xpath_user_email").String(),
		vUserEmail:   cfg.Section("ark").Key("input_xpath_v_user_email").String(),
		shipping:     cfg.Section("ark").Key("input_xpath_shipping").String(),
		payment:      cfg.Section("ark").Key("input_xpath_payment").String(),
		nextPage1:    cfg.Section("ark").Key("input_xpath_next_page1").String(),
		nextPage2:    cfg.Section("ark").Key("input_xpath_next_page2").String(),
		nextPage3:    cfg.Section("ark").Key("input_xpath_next_page3").String(),
		Tracer:       trace.New(os.Stdout),
	}
}

func (t *Ark) Run(user map[string]string) (err error) {

	// 在庫チェック→page.HTML()でやった方がいいのかな？
	if result := Check(t); result == false {
		return errors.New("在庫ないです")
	}

	// ブラウザ：chromeを指定して起動
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		fmt.Printf("Failed to start driver: %v", err)
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		fmt.Printf("Failed to open page: %v", err)
		return err
	}

	// cookieクリア
	if err := page.ClearCookies(); err != nil {
		fmt.Printf("Failed to clear cookies: %v", err)
		return err
	}
	sleep()

	// 商品ページに遷移
	if err := page.Navigate(t.targetUrl); err != nil {
		fmt.Printf("Failed to navigate: %v", err)
		return err
	}
	sleep()

	// カートに入れる、カート画面遷移
	if err := page.FindByClass(t.stockBtn).Submit(); err != nil {
		fmt.Printf("Failed to add to cart: %v", err)
		return err
	}
	sleep()
	if err := page.Navigate(t.addresseeUrl); err != nil {
		fmt.Printf("Failed to navigate: %v", err)
		return err
	}
	sleep()

	// 情報をばんばん入れてく
	// step1 宛先の入力
	if err := page.FindByXPath(t.name).Fill(user["name"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.nameKana).Fill(user["nameKana"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.zipcode1).Fill(user["zipcode1"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.zipcode2).Fill(user["zipcode2"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.pref).Select(user["pref"]); err != nil {
		fmt.Printf("Failed to select pref: %v", err)
		return err
	}
	if err := page.FindByXPath(t.city).Fill(user["city"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.street).Fill(user["street"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.building).Fill(user["building"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.phone).Fill(user["phone"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.userEmail).Fill(user["email"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	if err := page.FindByXPath(t.vUserEmail).Fill(user["email"]); err != nil {
		fmt.Printf("Failed to input: %v", err)
		return err
	}
	sleep()
	if err := page.FindByXPath(t.nextPage1).Click(); err != nil {
		fmt.Printf("Failed to submit at shipping form page: %v", err)
		return err
	}
	sleep()

	//step2 支払い方法・各種指定
	if err := page.FindByXPath(t.shipping).Click(); err != nil {
		fmt.Printf("Failed to select shipping method: %v", err)
		return err
	}
	if err := page.FindByXPath(t.payment).Click(); err != nil {
		fmt.Printf("Failed to select payment method: %v", err)
		return err
	}
	sleep()
	if err := page.FindByXPath(t.nextPage2).Click(); err != nil {
		fmt.Printf("Failed to submit at payment form page: %v", err)
		return err
	}
	sleep()

	//step3 注文確認画面→コメントアウト外しちゃうと買っちゃうはず、テストしてません。
	//if err := page.FindByXPath(t.nextPage3).Click(); err != nil {
	//	fmt.Printf("Failed to purchase: %v", err)
	//	return err
	//}
	sleep()

	// BOT判定の画像のやつがくるからうまくログインできない、できなくても買い物はできるので諦め
	//email := page.FindByXPath(t.inputMail)
	//sleep()
	//password := page.FindByXPath(t.inputPw)
	//sleep()
	//email.Fill(user["mailAddress"])
	//password.Fill(user["password"])
	//if err := page.FindByXPath(t.login).Submit();
	//	err != nil {
	//		fmt.Printf("Failed to login: %v", err)
	//		return err
	//}
	//sleep()

	return nil
}

func (t *Ark) getCheckInfo() map[string]string {
	return map[string]string{
		"targetUrl":  t.targetUrl,
		"checkPoint": t.stockBtn,
		"checkWord":  t.checkWord,
	}
}

func sleep() {
	time.Sleep(1 * time.Second)
	// 以下のようにするとエラー起きる、、、
	//time.Sleep(s * time.Second)
}
