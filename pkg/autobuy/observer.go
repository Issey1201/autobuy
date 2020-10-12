package autobuy

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// checkしたいのはArk型だけじゃないから、
// TargetSiteインターフェースを引数にして色々なサイトをポインタとして受け取りたい→interfaceをポインタで受け取るのはきつい？
//　→(t *Ark)check()bool　を、check(t TargetSite)boolとし、
//  tの中身はそれを取り出すメソッド(getCheckInfo)によって取得
func Check(t TargetSite) bool {

	info := t.getCheckInfo()
	// http.Getは構造体http.Responseのポインタを返してる
	// さらにres.Bodyは、io.ReadCloserインターフェース型である
	// io.ReadCloserインターフェース型はio.Readerとio.Closerインターフェース型を実装している
	//（インターフェースの中にインターフェースというのがしっくりこない、io.Readerとio.Closer両方の必須メソッドを持っているインターフェースということ？）
	// io.Readerインターフェースは右のメソッドを保持→func (Reader) Read(p []byte) (n int, err error)
	// io.Closerインターフェースは右のメソッドを保持→func (Closer) Close() error
	res, err := http.Get(info["targetUrl"])
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	// NewDocumentFromReader関数の引数は、io.Readerインターフェースを引数に取る
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stockStatus := doc.Find("." + info["checkPoint"]).Text()
	if strings.Index(stockStatus, info["checkWord"]) != -1 {
		// 在庫あり
		return true
	}
	// 在庫なし
	return false
}
