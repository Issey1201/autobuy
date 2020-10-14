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

type ark struct {
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
	User         user
	Tracer       trace.Tracer
}

type TargetSite interface {
	Run() error
	getCheckInfo() map[string]string
}

func NewArk(confPath string) *ark {

	var config struct {
		Ark  map[string]ark
		User user
	}

	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		fmt.Printf("Failed to open toml file: %v", err)
		return nil
	}

	return &ark{
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
		User:         config.User,
		Tracer:       trace.New(os.Stdout),
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
		if ret = page.Navigate(t.TargetUrl); ret != nil {
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
		if ret = page.FindByClass(t.StockBtn).Click(); ret != nil {
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
		if ret = page.Navigate(t.AddresseeUrl); ret != nil {
			return ret
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
		if ret = page.FindByXPath(t.Name).Fill(t.User.Name); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.NameKana).Fill(t.User.NameKana); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Zipcode1).Fill(t.User.Zipcode1); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Zipcode2).Fill(t.User.Zipcode2); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Pref).Select(t.User.Pref); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.City).Fill(t.User.City); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Street).Fill(t.User.Street); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Building).Fill(t.User.Building); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.Phone).Fill(t.User.Phone); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.UserEmail).Fill(t.User.Email); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.VUserEmail).Fill(t.User.Email); ret != nil {
			return ret
		}
		if ret = page.FindByXPath(t.NextPage1).Click(); ret != nil {
			return ret
		}
		return nil
	}, retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return time.Duration(n) * time.Second
	}), retry.Attempts(5)); err != nil {
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
	}), retry.Attempts(5)); err != nil {
		fmt.Printf("Failed at payment page: %v\n", err)
		return err
	}

	//step3 注文確認画面→コメントアウト外しちゃうと買っちゃうはず、テストしてません。
	//if err := page.FindByXPath(t.NextPage3).Click(); err != nil {
	//	fmt.Printf("Failed to purchase: %v", err)
	//	return err
	//}

	return nil
}

func (t *ark) getCheckInfo() map[string]string {
	return map[string]string{
		"targetUrl":  t.TargetUrl,
		"checkPoint": t.StockBtn,
		"checkWord":  t.CheckWord,
	}
}
