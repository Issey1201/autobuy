package autobuy

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Issey1201/pkg/trace"
	"github.com/avast/retry-go"
	"github.com/sclevine/agouti"
)

type Ark struct {
	Url          string `toml:"base_url"`
	TargetUrl    string `toml:"target_url"`
	AddresseeUrl string `toml:"addressee_url"`
	InputMail    string `toml:"email"`
	InputPw      string `toml:"password"`
	Login        string `toml:"login"`
	StockBtn     string `toml:"stock"`
	CheckWord    string `toml:"in_stock_word"`
	Name         string `toml:"name"`
	NameKana     string `toml:"name_kana"`
	Zipcode1     string `toml:"zipcode1"`
	Zipcode2     string `toml:"zipcode2"`
	Pref         string `toml:"pref"`
	City         string `toml:"city"`
	Street       string `toml:"street"`
	Building     string `toml:"building"`
	Phone        string `toml:"phone"`
	UserEmail    string `toml:"user_email"`
	VUserEmail   string `toml:"v_user_email"`
	Shipping     string `toml:"shipping"`
	Payment      string `toml:"payment"`
	NextPage1    string `toml:"next_page1"`
	NextPage2    string `toml:"next_page2"`
	NextPage3    string `toml:"next_page3"`
	Tracer       trace.Tracer
}

type TargetSite interface {
	Run(map[string]string) error
	getCheckInfo() map[string]string
}

type Config struct {
	//User map[string]User
	Ark map[string]Ark
}

func NewArk() *Ark {
	// 相対パスじゃなくて絶対パスにしたい
	// →main.goでもini.Load()をしていて、パスはmain.goに合わせてる感じ？
	// main.goでは相対パスなのに大して、こちらは絶対パス（？）少なくとも相対パスではない
	// 最初読み込んだら	var config Config２回目以降の読み込みは、１回目と同じ読み込み方で良いということ？
	var config Config
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		fmt.Printf("Failed to open toml file: %v", err)
		return nil
	}

	return &Ark{
		Url:          config.Ark["url"].Url,
		TargetUrl:    config.Ark["url"].TargetUrl,
		AddresseeUrl: config.Ark["url"].AddresseeUrl,
		InputMail:    config.Ark["xpath"].InputMail,
		InputPw:      config.Ark["xpath"].InputPw,
		Login:        config.Ark["xpath"].Login,
		StockBtn:     config.Ark["selector"].StockBtn,
		CheckWord:    config.Ark["other"].CheckWord,
		Name:         config.Ark["xpath"].Name,
		NameKana:     config.Ark["xpath"].NameKana,
		Zipcode1:     config.Ark["xpath"].Zipcode1,
		Zipcode2:     config.Ark["xpath"].Zipcode2,
		Pref:         config.Ark["xpath"].Pref,
		City:         config.Ark["xpath"].City,
		Street:       config.Ark["xpath"].Street,
		Building:     config.Ark["xpath"].Building,
		Phone:        config.Ark["xpath"].Phone,
		UserEmail:    config.Ark["xpath"].UserEmail,
		VUserEmail:   config.Ark["xpath"].VUserEmail,
		Shipping:     config.Ark["xpath"].Shipping,
		Payment:      config.Ark["xpath"].Payment,
		NextPage1:    config.Ark["xpath"].NextPage1,
		NextPage2:    config.Ark["xpath"].NextPage2,
		NextPage3:    config.Ark["xpath"].NextPage3,
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
		ret = page.Navigate(t.TargetUrl)
		return ret
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed to navigate: %v\n", err)
		return err
	}

	// カートに入れる、カート画面遷移
	if err = retry.Do(func() error {
		if ret = page.FindByClass(t.StockBtn).Submit(); ret != nil {
			return ret
		}
		if ret = page.Navigate(t.AddresseeUrl); ret != nil {
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
		if ret = page.FindByXPath(t.Name).Fill(user["name"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.NameKana).Fill(user["name_kana"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Zipcode1).Fill(user["zipcode1"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Zipcode2).Fill(user["zipcode2"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Pref).Select(user["pref"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.City).Fill(user["city"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Street).Fill(user["street"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Building).Fill(user["building"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Phone).Fill(user["phone"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.UserEmail).Fill(user["email"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.VUserEmail).Fill(user["email"]); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.NextPage1).Click(); ret != nil {
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
		if ret = page.FindByXPath(t.Shipping).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Payment).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.NextPage2).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(3)); err != nil {
		fmt.Printf("Failed at payment page: %v\n", err)
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
		"targetUrl":  t.TargetUrl,
		"checkPoint": t.StockBtn,
		"checkWord":  t.CheckWord,
	}
}
