package autobuy_test

import (
	"testing"

	"github.com/Issey1201/pkg/autobuy"
)

func TestYodobashi_Run(t *testing.T) {

	yodobashi := autobuy.NewYodobashi("./testdata/yodobashi.toml")

	tests := []struct {
		title string
		url   string
	}{
		{"InStock1", "https://www.yodobashi.com/product/100000001002677304/"},
		{"InStock2", "https://www.yodobashi.com/product/100000001002677304/"},
		{"OutOfStock", "https://www.yodobashi.com/product/100000001005857351/"},
		{"Valid", "https://www.yodobashi.com/product/99999999999999999/"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if err := yodobashi.Run(tt.url); err != nil {
				t.Fatalf("errorを返すべきでは無い: %v", err)
			}
		})
	}
}
