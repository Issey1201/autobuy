package autobuy

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CheckResponse struct {
	StockStatus bool
	Url         string
}

func CheckStock(t TargetSite, targetUrl string, ch chan CheckResponse, wg *sync.WaitGroup) {
	for {
		result := Check(t, targetUrl)
		cr := &CheckResponse{
			StockStatus: result,
			Url:         targetUrl,
		}
		// １つ以上あるCheckStockゴルーチンの中で、１つでも在庫があったのであれば、
		// 他のゴルーチン終了させてチャネル送信したくない
		ch <- *cr

		if cr.StockStatus == true {
			// 現状チャネルをcloseする時はどこか１つのURLにて在庫があった場合であるが、
			// 上記以外のURLに対するゴルーチンでチャネルを送信してしまうのでエラーを吐く場合がある
			// 在庫ないURLを全て最初に処理され、最後に在庫あったURLのゴルーチンが処理された場合に限りエラーを吐かない...1'
			// 1'の場合、在庫がないURLのゴルーチンが勝手に終了されているのは謎、、、
			// 在庫がないURLのゴルーチンのbreakされてないからfor文はまだ回っているはず
			close(ch)
			//wg.Done()
			_ = wg
			break
		} else {
			time.Sleep(10 * time.Second)
		}
	}
	//wg.Wait()
	//defer close(f)
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
