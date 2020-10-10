package autobuy

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Check() string {
	url := "https://www.ark-pc.co.jp/i/20106260/"
	//url := "https://ja.wikipedia.org/wiki/SCADA"

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	selection := doc.Find(".btn-addcart")
	return selection.Text()
}