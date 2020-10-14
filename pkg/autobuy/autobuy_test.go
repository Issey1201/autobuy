package autobuy_test

import (
	"testing"

	"github.com/Issey1201/pkg/autobuy"
)

func TestArk_Run(t *testing.T) {

	ark := autobuy.NewArk("./testdata/ark.toml")

	if err := ark.Run(); err != nil {
		t.Errorf("errorを返すべきでは無い: %v", err)
	}
}
