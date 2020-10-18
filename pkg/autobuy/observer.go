package autobuy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CheckResponse struct {
	StockStatus bool
	Url         string
}

func Check(t TargetSite, targetUrl string, ch chan CheckResponse, done chan struct{}) {
	tick := time.NewTicker(1 * time.Minute)
	defer tick.Stop()
	for {
		select {
		case <-done:
			return
		case <-tick.C:
			result := CheckStock(t, targetUrl)
			cr := &CheckResponse{
				StockStatus: result,
				Url:         targetUrl,
			}
			ch <- *cr
		}
	}
}

func CheckStock(t TargetSite, targetUrl string) bool {

	info := t.getCheckInfo()
	res, err := http.Get(targetUrl)
	if err != nil {
		fmt.Printf("failed to get html: %v", err)
		return false
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
