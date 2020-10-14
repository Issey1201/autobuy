package autobuy

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Issey1201/pkg/trace"
	"github.com/avast/retry-go"
	"github.com/sclevine/agouti"
)

type arkUrl struct {
	BaseUrl      string `toml:"base_url"`
	TargetUrl    string `toml:"target_url"`
	AddresseeUrl string `toml:"addressee_url"`
}

type arlSelector struct {
	StockBtn string `toml:"stock"`
}

type arkXpath struct {
	InputMail  string `toml:"email"`
	InputPw    string `toml:"password"`
	Login      string `toml:"login"`
	Name       string `toml:"name"`
	NameKana   string `toml:"name_kana"`
	Zipcode1   string `toml:"zipcode1"`
	Zipcode2   string `toml:"zipcode2"`
	Pref       string `toml:"pref"`
	City       string `toml:"city"`
	Street     string `toml:"street"`
	Building   string `toml:"building"`
	Phone      string `toml:"phone"`
	UserEmail  string `toml:"user_email"`
	VUserEmail string `toml:"v_user_email"`
	Shipping   string `toml:"shipping"`
	Payment    string `toml:"payment"`
	NextPage1  string `toml:"next_page1"`
	NextPage2  string `toml:"next_page2"`
	NextPage3  string `toml:"next_page3"`
}

type arkOther struct {
	CheckWord string `toml:"in_stock_word"`
}

type arkConf struct {
	Url      arkUrl
	Selector arlSelector
	Xpath    arkXpath
	Other    arkOther
	User     user
}

type ark struct {
	Config arkConf
	Tracer trace.Tracer
}

func NewArk(confPath string) *ark {

	var config arkConf
	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		fmt.Printf("Failed to open toml file: %v", err)
		return nil
	}

	return &ark{
		Config: config,
		Tracer: trace.New(os.Stdout),
	}
}

func (t *ark) Run() (err error) {

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
		if ret = page.Navigate(t.Config.Url.TargetUrl); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// カートに入れる
	if err = retry.Do(func() error {
		if ret = page.FindByClass(t.Config.Selector.StockBtn).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// カート画面遷移
	if err = retry.Do(func() error {
		if ret = page.Navigate(t.Config.Url.AddresseeUrl); ret != nil {
			return ret
		}
		if url, _ := page.URL(); url != t.Config.Url.AddresseeUrl {
			return errors.New("failed to move page")
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// step1 宛先の入力
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.Name).Fill(t.Config.User.Name); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.NameKana).Fill(t.Config.User.NameKana); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Zipcode1).Fill(t.Config.User.Zipcode1); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Zipcode2).Fill(t.Config.User.Zipcode2); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Pref).Select(t.Config.User.Pref); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.City).Fill(t.Config.User.City); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Street).Fill(t.Config.User.Street); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Building).Fill(t.Config.User.Building); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Phone).Fill(t.Config.User.Phone); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.UserEmail).Fill(t.Config.User.Email); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.VUserEmail).Fill(t.Config.User.Email); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.NextPage1).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed at input page: %v\n", err)
		return err
	}

	// step2 支払い方法・各種指定
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.Shipping).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.Payment).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.NextPage2).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed at payment page: %v\n", err)
		return err
	}

	//// step3 注文確認画面→コメントアウト外しちゃうと買っちゃうはず、テストしてません。
	//if err := page.FindByXPath(t.Config.Xpath.NextPage3).Click(); err != nil {
	//	fmt.Printf("Failed to purchase: %v", err)
	//	return err
	//}

	return nil
}

func (t *ark) getCheckInfo() map[string]string {
	return map[string]string{
		"targetUrl":  t.Config.Url.TargetUrl,
		"checkPoint": t.Config.Selector.StockBtn,
		"checkWord":  t.Config.Other.CheckWord,
	}
}
