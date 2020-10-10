package autobuy

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Ark) Check() bool {

	// http.Getは構造体http.Responseのポインタを返してる
	// さらにres.Bodyは、io.ReadCloserインターフェース型である
	// io.ReadCloserインターフェース型はio.Readerとio.Closerインターフェース型を実装している
	//（インターフェースの中にインターフェースというのがしっくりこない、io.Readerとio.Closer両方の必須メソッドを持っているインターフェースということ？）
	// io.Readerインターフェースは右のメソッドを保持→func (Reader) Read(p []byte) (n int, err error)
	// io.Closerインターフェースは右のメソッドを保持→func (Closer) Close() error
	res, err := http.Get(t.targetUrl)
	if err != nil {
		panic(err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	// NewDocumentFromReader関数の引数は、io.Readerインターフェースを引数に取る
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stockStatus := doc.Find(t.checkPoint).Text()
	if strings.Index(stockStatus, t.checkWord) != -1 {
		// 在庫なし
		return false
	}
	// 在庫あり
	return true
}