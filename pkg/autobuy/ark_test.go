package autobuy_test

import (
	"testing"

	"github.com/Issey1201/pkg/autobuy"
)

func TestArk_Run(t *testing.T) {

	ark := autobuy.NewArk("./testdata/ark.toml")

	tests := []struct {
		title string
		url   string
	}{
		{"InStock", "https://www.ark-pc.co.jp/i/50293029/"},
		{"OutOfStock", "https://www.ark-pc.co.jp/i/20106274/"},
		{"Valid", "https://www.ark-pc.co.jp/i/99999999999/"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if err := ark.Run(tt.url); err != nil {
				t.Fatalf("errorを返すべきでは無い: %v", err)
			}
		})
	}
}
