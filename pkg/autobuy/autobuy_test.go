package autobuy

import (
	"testing"
)

func TestArk_Run(t *testing.T) {

	user := map[string]string{
		"password": "sampLe_123",
		"name":     "仮名",
		"nameKana": "カリメイ",
		"zipcode1": "121",
		"zipcode2": "1234",
		"pref":     "東京都",
		"city":     "足立区○○町",
		"street":   "1丁目-11-22",
		"building": "○○ビル7Ｆ",
		"phone":    "0352987020",
		"email":    "ark@sample.co.jp",
	}

	// ログイン情報とURLは引数とするべきか？それとも構造体に格納すべきか？
	// arkに関する固定情報は構造体、user情報などarkでも情報が状況により変わってくるのは引数が良い？
	ark := NewArk()
	if err := ark.Run(user); err != nil {
		t.Errorf("errorを返すべきでは無い: %v", err)
	}
}
