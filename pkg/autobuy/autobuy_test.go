package autobuy

import (
	"testing"
)

func TestArk_Run(t *testing.T) {

	// ログイン情報とURLは引数とするべきか？それとも構造体に格納すべきか？
	// arkに関する固定情報は構造体、user情報などarkでも情報が状況により変わってくるのは引数が良い？
	ark := NewArk()
	if err := ark.Run(); err != nil {
		t.Errorf("errorを返すべきでは無い: %v", err)
	}
}
