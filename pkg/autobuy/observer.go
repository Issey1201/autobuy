package autobuy

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func CheckStock(t TargetSite, targetUrl string, resultCh chan<- bool) {
	for {
		if result := Check(t, targetUrl); result == false {
			resultCh <- result
			time.Sleep(1 * time.Minute)
		} else {
			resultCh <- result
			break
		}
	}
	close(resultCh)
}

// checkしたいのはArk型だけじゃないから、
// TargetSiteインターフェースを引数にして色々なサイトをポインタとして受け取りたい→interfaceをポインタで受け取るのはきつい？
//　→(t *Ark)check()bool　を、check(t TargetSite)boolとし、
//  tの中身はそれを取り出すメソッド(getCheckInfo)によって取得
// 引数の型で t *TargetSite はなぜダメ？
func Check(t TargetSite, targetUrl string) bool {

	info := t.getCheckInfo()
	res, err := http.Get(targetUrl)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("failed to read a new document: %v", err)
	}

	stockStatus := doc.Find("." + info["checkPoint"]).Text()
	if strings.Index(stockStatus, info["checkWord"]) != -1 {
		// 在庫あり
		return true
	}
	// 在庫なし
	return false
}
