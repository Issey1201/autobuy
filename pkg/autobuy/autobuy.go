package autobuy

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Issey1201/pkg/trace"
	"github.com/avast/retry-go"
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

	var ret error

	// 商品ページに遷移
	if err = retry.Do(func() error {
		ret = page.Navigate(t.targetUrl)
		return ret
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed to navigate: %v\n", err)
		return err
	}

	// カートに入れる、カート画面遷移
	if err = retry.Do(func() error {
		if ret = page.FindByClass(t.stockBtn).Submit(); ret != nil {
			return ret
		}
		if ret = page.Navigate(t.addresseeUrl); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// step1 宛先の入力
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.name).Fill(user["name"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.nameKana).Fill(user["nameKana"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.zipcode1).Fill(user["zipcode1"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.zipcode2).Fill(user["zipcode2"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.pref).Select(user["pref"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.city).Fill(user["city"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.street).Fill(user["street"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.building).Fill(user["building"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.phone).Fill(user["phone"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.userEmail).Fill(user["email"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.vUserEmail).Fill(user["email"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.nextPage1).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed at input page: %v\n", err)
		return err
	}

	//step2 支払い方法・各種指定
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.shipping).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.payment).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.nextPage2).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed at input page: %v\n", err)
		return err
	}

	//step3 注文確認画面→コメントアウト外しちゃうと買っちゃうはず、テストしてません。
	//if err := page.FindByXPath(t.nextPage3).Click(); err != nil {
	//	fmt.Printf("Failed to purchase: %v", err)
	//	return err
	//}

	return nil
}

func (t *Ark) getCheckInfo() map[string]string {
	return map[string]string{
		"targetUrl":  t.targetUrl,
		"checkPoint": t.stockBtn,
		"checkWord":  t.checkWord,
	}
}
