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

type yodobashiUrl struct {
	BaseUrl   string   `toml:"base_url"`
	TargetUrl []string `toml:"target_url"`
}

type yodobashiSelector struct {
	Stock string `toml:"stock"`
}

type yodobashiXpath struct {
	Add2Cart     string `toml:"add2cart"`
	PurchasePage string `toml:"purchase_page"`
	InputPw      string `toml:"password"`
	Email        string `toml:"email"`
	Login        string `toml:"login"`
	Payment      string `toml:"payment"`
	NextPage1    string `toml:"next_page1"`
	NextPage2    string `toml:"next_page2"`
	NextPage3    string `toml:"next_page3"`
	OrderConfirm string `toml:"order_confirm"`
}

type yodobashiOther struct {
	CheckWord string `toml:"in_stock_word"`
}

type yodobashiConf struct {
	Url      yodobashiUrl
	Selector yodobashiSelector
	Xpath    yodobashiXpath
	Other    yodobashiOther
	User     user
}

type yodobashi struct {
	Config yodobashiConf
	Tracer trace.Tracer
}

func NewYodobashi(confPath string) *yodobashi {

	var config yodobashiConf
	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		fmt.Printf("Failed to open toml file: %v", err)
		return nil
	}

	return &yodobashi{
		Config: config,
		Tracer: trace.New(os.Stdout),
	}
}

func (t *yodobashi) Run(targetUrl string) (err error) {

	attempt := uint(10)

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
		if ret = page.Navigate(targetUrl); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// カートに入れる
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.Add2Cart).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// 購入手続きに進む(1)
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.PurchasePage).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// 次へ進む
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.NextPage1).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed to add to cart: %v\n", err)
		return err
	}

	// ログイン
	if err = retry.Do(func() error {
		if ret = page.FindByXPath(t.Config.Xpath.Email).Fill(t.Config.User.Email); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.InputPw).Fill(t.Config.User.Password); ret != nil {
			return ret
		}
		time.Sleep(5 * time.Second)

		if ret = page.FindByXPath(t.Config.Xpath.Login).Click(); ret != nil {
			return ret
		}
		time.Sleep(5 * time.Second)
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed at input page: %v\n", err)
		return err
	}

	// 購入手続きに進む(2)
	if err = retry.Do(func() error {
		time.Sleep(5 * time.Second)
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

	// 支払い方法の選択
	if err = retry.Do(func() error {
		time.Sleep(5 * time.Second)
		if ret = page.FindByXPath(t.Config.Xpath.Payment).Click(); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Config.Xpath.NextPage3).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(attempt)); err != nil {
		fmt.Printf("Failed at payment page: %v\n", err)
		return err
	}

	// step3 注文確定
	if os.Getenv("ENV") == "prod" {
		if err := page.FindByXPath(t.Config.Xpath.OrderConfirm).Click(); err != nil {
			fmt.Printf("Failed to purchase: %v", err)
			return err
		}
	}

	return nil
}

func (t *yodobashi) getCheckInfo() map[string]string {
	return map[string]string{
		"checkPoint": t.Config.Selector.Stock,
		"checkWord":  t.Config.Other.CheckWord,
	}
}
