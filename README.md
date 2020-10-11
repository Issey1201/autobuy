# autobuy
自動購入スクリプト

### quick start
1.golangのinstall
  https://golang.org/doc/install

2.以下のpackageをinstall
	"github.com/go-ini/ini"
	"github.com/sclevine/agouti"
  "github.com/PuerkitoBio/goquery"
  "github.com/bwmarrin/discordgo"

3.chromedriverをinstall
windowsの人は以下を参考に
https://www.mittsu-kosen.com/chromedriver%E3%82%92windows10%E3%81%A7%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB%E3%81%99%E3%82%8B%E6%96%B9%E6%B3%95%E3%80%90%E7%94%BB%E5%83%8F%E4%BB%98%E3%81%8D%E3%80%91/
コマンドプロンプトとかでどのディレクトリでもいいので[chromedriver]と叩いてchromedriverが反応するか確認してください。もし、
「command not found: chromedriver」 と出るのであれば、PCの再起動なり、pathが通っているかを確認してください。[chromedriver path windows]とかで検索。

4.config.iniの以下の情報を自分のやつに変更してください
# user_info
user_email = ark@sample.co.jp
user_password =
user_name = 嗚呼句　太朗
user_name_kana = アアク　タロウ
user_zipcode1 = 123
user_zipcode2 = 0851
user_pref = 東京都
user_city = 足立区○○町
user_street = 1丁目-11-22
user_building = ○○ビル7Ｆ
user_phone = 0352987020
# url #target_url→欲しい商品が記載されているURL
target_url = https://www.ark-pc.co.jp/i/11501894/

5.autobuyのディレクトリで以下のコマンドを実行
go run main.gp
